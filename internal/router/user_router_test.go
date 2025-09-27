package router

import (
	"testing"
)

func TestFind(t *testing.T) {
	tester := NewRouterTester(t)

	t.Run("Should return the user with a 200 status code", func(t *testing.T) {
		tester.
			Get("/user/1").
			AssertStatus(200).
			AssertBodyKeyIsString("id").
			AssertBodyKeyIsString("name").
			AssertBodyKeyIsString("username").
			AssertBodyKeyIsDate("created_at").
			AssertBodyKeyIsDate("updated_at")
	})
}
