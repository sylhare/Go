package internal

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestContainer(t *testing.T) {
	intContainer := Container[int]{value: 42}
	assert.Equal(t, 42, intContainer.GetValue())

	stringContainer := Container[string]{value: "hello"}
	assert.Equal(t, "hello", stringContainer.GetValue())
}

func TestFirstElement(t *testing.T) {
	intSlice := []int{1, 2, 3}
	assert.Equal(t, 1, FirstElement(intSlice))

	stringSlice := []string{"a", "b", "c"}
	assert.Equal(t, "a", FirstElement(stringSlice))
}

func TestSwap(t *testing.T) {
	a, b := 1, 2
	Swap(&a, &b)
	assert.Equal(t, 2, a)
	assert.Equal(t, 1, b)

	x, y := "hello", "world"
	Swap(&x, &y)
	assert.Equal(t, "world", x)
	assert.Equal(t, "hello", y)
}

func TestPrintValueWithStruct(t *testing.T) {
	r := RegularStruct{name: "test"}

	result := StringifyValue(r, 42)
	assert.Equal(t, "test: Value: 42", result)

	result = StringifyValue(r, "hello")
	assert.Equal(t, "test: Value: hello", result)
}

func TestPtr(t *testing.T) {
	testCases := []struct {
		name  string
		value any
	}{
		{"int", 42},
		{"string", "hello"},
		{"struct", RegularStruct{name: "test"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ptrValue := ptr(tc.value)
			assert.Equal(t, &tc.value, ptrValue)
			assert.Equal(t, tc.value, *ptrValue)
			assert.Equal(t, reflect.Ptr, reflect.TypeOf(ptrValue).Kind())
		})
	}
}
