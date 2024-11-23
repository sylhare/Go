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
func Modify(a *int) int {
	*a += 1
	return *a
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
