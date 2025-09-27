package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/omatheuscaetano/planus-api/internal/dependency_container"
	"github.com/omatheuscaetano/planus-api/pkg/env"
	"github.com/stretchr/testify/assert"
)

type RouterTester struct {
	t         *testing.T
	engine    *gin.Engine
	response  *httptest.ResponseRecorder
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
	tester.response = httptest.NewRecorder()
	req, _ := http.NewRequest("GET", path, nil)
	tester.engine.ServeHTTP(tester.response, req)
	return tester
}

func (tester *RouterTester) AssertStatus(expected int) *RouterTester {
	assert.Equal(tester.t, expected, tester.response.Code)
	return tester
}

