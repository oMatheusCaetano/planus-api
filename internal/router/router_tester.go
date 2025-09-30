package router

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/omatheuscaetano/planus-api/internal/dependency_container"
	mongo "github.com/omatheuscaetano/planus-api/pkg/db"
	"github.com/omatheuscaetano/planus-api/pkg/env"
	"github.com/stretchr/testify/assert"
)

type RouterTester struct {
	t            *testing.T
	engine       *gin.Engine
	response     *httptest.ResponseRecorder
	responseBody map[string]any
	mongo        *mongo.MongoDB
}

func NewRouterTester(t *testing.T) *RouterTester {
	env.Load()

	tester := &RouterTester{
		t:      t,
		engine: gin.New(),
		mongo:  mongo.NewTestMongoDb(),
	}

	err := tester.mongo.Connect(context.Background())
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	instances := dependency_container.InstantiateAll(&dependency_container.InstancesConfig{
		Mongo: tester.mongo,
	})
	AllRoutes(tester.engine, instances)

	return tester
}

func (tester *RouterTester) InsertToTable(table string, payload interface{}) any {
	res, err := tester.mongo.Collection(table).InsertOne(context.Background(), payload)
	if err != nil {
		tester.t.Fatalf("failed to create user: %v", err)
	}

	if res.InsertedID == nil {
		tester.t.Fatalf("failed to create user: no ID returned")
	}

	return res.InsertedID
}

func (tester *RouterTester) DropTable(table string) *RouterTester {
	err := tester.mongo.Collection(table).Drop(context.Background())
	if err != nil {
		tester.t.Fatalf("failed to drop %s collection: %v", table, err)
	}
	return tester
}

func (tester *RouterTester) FindFromTable(table string, filter interface{}) map[string]any {
	row := tester.mongo.Collection(table).FindOne(context.Background(), filter)

	if row.Err() != nil {
		tester.t.Fatalf("failed to find document in %s collection: %v", table, row.Err())
	}

	var entity map[string]any
	err := row.Decode(&entity)
	if err != nil {
		tester.t.Fatalf("failed to decode document from %s collection: %v", table, err)
	}
	return entity
}

func (tester *RouterTester) Get(path string) *RouterTester {
	tester.responseBody = nil
	tester.response = httptest.NewRecorder()
	req, _ := http.NewRequest("GET", path, nil)
	tester.engine.ServeHTTP(tester.response, req)
	return tester
}

func (tester *RouterTester) Delete(path string) *RouterTester {
	tester.responseBody = nil
	tester.response = httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", path, nil)
	tester.engine.ServeHTTP(tester.response, req)
	return tester
}

func (tester *RouterTester) Post(path string, body any) *RouterTester {
	tester.responseBody = nil
	tester.response = httptest.NewRecorder()
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", path, strings.NewReader(string(jsonBody)))
	req.Header.Set("Content-Type", "application/json")

	tester.engine.ServeHTTP(tester.response, req)
	return tester
}

func (tester *RouterTester) Body() *map[string]any {
	if (tester.responseBody == nil) {
		bodyString := tester.response.Body.String()
		bodyMap := make(map[string]any)
		json.Unmarshal([]byte(bodyString), &bodyMap)
		tester.responseBody = bodyMap
	}
	return &tester.responseBody
}

func (tester *RouterTester) AssertEqual(expected interface{}, actual interface{}) *RouterTester {
	assert.Equal(tester.t, expected, actual)
	return tester
}

func (tester *RouterTester) AssertStatus(expected int) *RouterTester {
	assert.Equal(tester.t, expected, tester.response.Code)
	return tester
}

func (tester *RouterTester) AssertBodyContains(expected string) *RouterTester {
	assert.Contains(tester.t, tester.Body(), expected)
	return tester
}

func (tester *RouterTester) AssertBodyKeyIsString(expected string) *RouterTester {
	body := *tester.Body()
	assert.Contains(tester.t, body, expected)
	assert.IsType(tester.t, "", body[expected])
	return tester
}

func (tester *RouterTester) AssertBodyKeyIsNumeric(expected string) *RouterTester {
	body := *tester.Body()
	assert.Contains(tester.t, body, expected)
	assert.IsType(tester.t, float64(0), body[expected])
	return tester
}

func (tester *RouterTester) AssertBodyValue(key string, value any) *RouterTester {
	body := *tester.Body()

	val, err := getNestedValue(body, key)
	if err != nil {
		tester.t.Errorf("failed to get value for key '%s': %v", key, err)
		return tester
	}

	assert.Equal(tester.t, value, val)
	return tester
}

func (tester *RouterTester) AssertBodyValueIsDateTime(key string, value time.Time) *RouterTester {
	body := *tester.Body()

	val, err := getNestedValue(body, key)
	if err != nil {
		tester.t.Errorf("failed to get value for key '%s': %v", key, err)
		return tester
	}

	assert.IsType(tester.t, "", val)

	parsedTime, err := time.Parse(time.RFC3339, val.(string))
	assert.NoError(tester.t, err)
	assert.Equal(tester.t, value.Format("2006-01-02 15:04:05"), parsedTime.Format("2006-01-02 15:04:05"))

	return tester
}

func (tester *RouterTester) AssertBodyKeyIsJWT(key string) *RouterTester {
	body := *tester.Body()

	val, err := getNestedValue(body, key)
	if err != nil {
		tester.t.Errorf("failed to get value for key '%s': %v", key, err)
		return tester
	}

	assert.IsType(tester.t, "", val)
	tokenParts := strings.Split(val.(string), ".")
	assert.Equal(tester.t, 3, len(tokenParts), "Invalid JWT format")

	return tester
}

func getNestedValue(data map[string]any, path string) (any, error) {
	keys := strings.Split(path, ".")
	var current any = data

	for i, k := range keys {
		currentMap, ok := current.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("expected map at '%s', got %T", strings.Join(keys[:i], "."), current)
		}

		val, exists := currentMap[k]
		if !exists {
			return nil, fmt.Errorf("key '%s' not found at path '%s'", k, strings.Join(keys[:i], "."))
		}

		if i == len(keys)-1 {
			return val, nil
		}
		current = val
	}

	return nil, fmt.Errorf("invalid path '%s'", path)
}
