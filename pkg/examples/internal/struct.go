package internal

type Counter struct {
	Value int
}

// Increment with value receiver (does not modify the actual value)
func (c Counter) Increment() int {
	c.Value++
	return c.Value
}

// Decrement with pointer receiver (modifies the actual value)
func (c *Counter) Decrement() {
	c.Value--
}
