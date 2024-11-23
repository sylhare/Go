package internal

import (
	"context"
	"fmt"
	"sync"
)

func tutorial() {
	defer fmt.Println(", world!")

	const Pi = 3.14159
	fmt.Println(Pi)

	var a int
	a = 1
	fmt.Println(a)

	b := 2
	fmt.Println(b)
	hello, _ := example()
	fmt.Printf(hello)
}
func example() (hello string, error error) {
	if err := doSomething(); err != nil {
		return
	}
	return "hello", nil
}
