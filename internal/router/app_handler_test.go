package router

import (
	"testing"
)

func TestWelcome(t *testing.T) {
	tester := NewRouterTester(t)

	t.Run("should return a welcome message", func(t *testing.T) {
		tester.Get("/").AssertStatus(200)
	})
}
