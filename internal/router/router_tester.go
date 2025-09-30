package router

import (
	"context"
	"encoding/json"
	"log"
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
		log.Println("Response Body:", bodyString) // Debugging line
		bodyMap := make(map[string]any)
		json.Unmarshal([]byte(bodyString), &bodyMap)
		tester.responseBody = bodyMap
	}
	return &tester.responseBody
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

func (tester *RouterTester) AssertBodyKeyIsDate(expected string) *RouterTester {
	body := *tester.Body()
	assert.Contains(tester.t, body, expected)
	assert.IsType(tester.t, "", body[expected])

	_, err := time.Parse(time.RFC3339, body[expected].(string))
	assert.NoError(tester.t, err)

	return tester
}


func (tester *RouterTester) AssertBodyValue(key string, value any) *RouterTester {
	body := *tester.Body()

	// Split the key by dots to handle nested paths
	keys := strings.Split(key, ".")

	// Navigate through the nested structure
	var current any = body
	for i, k := range keys {
		// Check if current is a map
		currentMap, ok := current.(map[string]any)
		if !ok {
			tester.t.Errorf("Expected map at path '%s', but got %T", strings.Join(keys[:i], "."), current)
			return tester
		}

		// Check if key exists in current map
		assert.Contains(tester.t, currentMap, k)

		// If this is the last key, compare the value
		if i == len(keys)-1 {
			assert.Equal(tester.t, value, currentMap[k])
		} else {
			// Otherwise, move to the next level
			current = currentMap[k]
		}
	}

	return tester
}

func (tester *RouterTester) AssertBodyValueIsDateTime(key string, value time.Time) *RouterTester {
	body := *tester.Body()
	assert.Contains(tester.t, body, key)
	assert.IsType(tester.t, "", body[key])

	parsedTime, err := time.Parse(time.RFC3339, body[key].(string))
	assert.NoError(tester.t, err)
	assert.Equal(tester.t, value.Format("2006-01-02 15:04:05"), parsedTime.Format("2006-01-02 15:04:05"))

	return tester
}
