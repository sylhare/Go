package internal

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCounter(t *testing.T) {
	c := Counter{Value: 10}

	// Increment using value receiver
	var incremented = c.Increment()
	fmt.Println("After Increment:", c.Value) // Output: 10 - didn't change
	assert.Equal(t, 10, c.Value)
	assert.Equal(t, 11, incremented)

	// Decrement using pointer receiver
	c.Decrement()
	fmt.Println("After Decrement:", c.Value) // Output: 9 - changed
	assert.Equal(t, 9, c.Value)
}
