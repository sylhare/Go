package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRunTutorial(t *testing.T) {
	tutorial()
}

func TestRunGoroutines(t *testing.T) {
	goroutines()
}

func TestModify(t *testing.T) {
	a, b := 2, 2
	t.Run("Test Modify", func(t *testing.T) {
		assert.Equal(t, 3, Modify(&a))
		assert.Equal(t, 3, a)
	})

	t.Run("Test ModifyForReal", func(t *testing.T) {
		assert.Equal(t, 3, ModifyForReal(b))
		assert.Equal(t, 2, b)
	})

	t.Run("Test ModifyBoth", func(t *testing.T) {
		modifiedA, modifiedB := ModifyBoth(&a, b)
		assert.Equal(t, 3, modifiedA)
		assert.Equal(t, 3, modifiedB)
		assert.Equal(t, 3, a)
		assert.Equal(t, 2, b)
	})
}

// TestUserHandlerImplementsHandler won't compile if UserHandler doesn't implement Handler
func TestUserHandlerImplementsHandler(t *testing.T) {
	var _ Handler = UserHandler{}
}