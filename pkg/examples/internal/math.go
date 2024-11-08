package internal

import "fmt"

func add(a int, b int) int {
	return a + b
}

func Numbers() {
	lang, ten := "golang", 10

	fmt.Printf("\nlet's try some %v! \n", lang)
	fmt.Printf("%d = %b + %d? \n", ten, 5, 5)
	fmt.Printf("%d = %X / %X? \n", ten, 50, 5)

	fmt.Println("\nlet's loop!")
	for i := 1; i < ten; i++ {
		fmt.Printf("%d \t %b \t %X \n", i, i, i)
	}
}
