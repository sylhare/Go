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

	p := Person{Name: "Alice", Age: 30}
	fmt.Println(p.Name)

	s := []int{1, 2, 3}
	fmt.Println(s[0])
	s = append(s, 4)
	fmt.Println(s)

	m := map[string]int{
		"one": 1,
		"two": 2,
	}
	fmt.Println(m["one"])

	m["three"] = 3
	fmt.Println(m)

	hello, _ := example()
	fmt.Printf(hello)
}

func doSomething() error {
	return nil
	//return error(fmt.Errorf("Raised an error"))
}

func Modify(a *int) int {
	*a += 1
	return *a // dereference the pointer to get the value
}

func ModifyForReal(a int) int {
	a += 1
	return a
}

func ModifyBoth(a *int, b int) (int, int) {
	*a += 1
	b += 1
	return *a, b
}

func example() (hello string, error error) {
	if err := doSomething(); err != nil {
		return
	}
	return "hello", nil
}

// Structs

type Person struct {
	Name string
	Age  int
}

// Interface

type Handler interface {
	Handle()
}

type UserHandler struct{}

func (h UserHandler) Handle() {
	fmt.Println("hello")
}

// Goroutines

func goroutines() {
	fmt.Println("hello 0")

	go fmt.Println("hello 1")

	// Context
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-ctx.Done()
		fmt.Println("~~~> done")
	}()
	cancel()

	// Sync
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("hello 2")
	}()

	wg.Wait()

	// Channel
	done := make(chan bool)

	go func() {
		fmt.Println("hello 3")
		done <- true
	}()

	<-done
	ch := make(chan int)
	go func() { ch <- 42 }()
	fmt.Println(<-ch)

	// Select
	ch1 := make(chan int)
	ch2 := make(chan int)
	go func() { ch1 <- 43 }()
	go func() { ch2 <- 44 }()
	select {
	case v1 := <-ch1:
		fmt.Println(v1)
	case v2 := <-ch2:
		fmt.Println(v2)
	}

	// will get the second value
	select {
	case v1 := <-ch1:
		fmt.Println(v1)
	case v2 := <-ch2:
		fmt.Println(v2)
	}
}

func exploreContext() error {
	const valueKey = "valueKey"
	const functionKey = "functionKey"
	myFunc := func() {
		fmt.Println("Function called from context")
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, valueKey, "myValue")
	ctx = context.WithValue(ctx, functionKey, myFunc)

	value := ctx.Value(valueKey)
	if value == nil {
		return fmt.Errorf("value not found in context")
	}
	fmt.Println("Value found in context:", value)
	function := ctx.Value(functionKey)
	if function == nil {
		return fmt.Errorf("function not found in context")
	}
	function.(func())()

	return nil
}
