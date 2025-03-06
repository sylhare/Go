package internal

import "fmt"

type Container[T any] struct {
	value T
}

func (c *Container[T]) GetValue() T {
	return c.value
}

func (c *Container[T]) SetValue(value T) {
	c.value = value
}

func FirstElement[T any](slice []T) T {
	return slice[0]
}

func Swap[T any](a, b *T) {
	*a, *b = *b, *a
}

type RegularStruct struct {
	name string
}

// StringifyValue is not in RegularStruct because methods on structs cannot have their own type parameters
// func (r *RegularStruct) PrintValue[T any](value T) string { // does not compile
func StringifyValue[T any](r RegularStruct, value T) string {
	return fmt.Sprintf("%s: Value: %v", r.name, value)
}

func ptr[T any](v T) *T {
	return &v
}

type Empty interface{}

type A interface {
	One()
	Two()
	Values() A
}

type StructA struct{ values []string }

func (s *StructA) One()      {}
func (s *StructA) Two()      {}
func (s *StructA) Values() A { return s }

type B interface {
	One()
	Values() B
}

type StructB struct{ values []int }

func (s *StructB) One()      {}
func (s *StructB) Values() B { return s }

type C[T any] interface {
	One()
	Values() T
}
