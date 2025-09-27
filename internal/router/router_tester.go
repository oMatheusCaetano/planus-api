package router

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/omatheuscaetano/planus-api/internal/dependency_container"
	"github.com/omatheuscaetano/planus-api/pkg/env"
	"github.com/stretchr/testify/assert"
)

type RouterTester struct {
	t            *testing.T
	engine       *gin.Engine
	response     *httptest.ResponseRecorder
	responseBody map[string]any
}

func NewRouterTester(t *testing.T) *RouterTester {
	env.Load()

	tester := &RouterTester{
		t:      t,
		engine: gin.New(),
	}

	instances := dependency_container.InstantiateAll()
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

func (tester *RouterTester) Body() *map[string]any {
	if (tester.responseBody == nil) {
		bodyString := tester.response.Body.String()
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
	assert.Contains(tester.t, body, key)
	assert.Equal(tester.t, value, body[key])
	return tester
}
