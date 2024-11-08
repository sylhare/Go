package internal

import "fmt"

var integer int
var word string
var float float64
var boolean bool

var WORLD = "world"

func Variables() {
	ten, lang := 37.5, "golang"
	var a, b, c = 1, 2, true

	fmt.Printf("%v of type %T \n", ten, ten)
	fmt.Printf("%v of type %T \n", lang, lang)
	fmt.Println()

	fmt.Printf("Here is an empty: \n- integer %v \n", integer)
	fmt.Printf("- string '%v' \n", word)
	fmt.Printf("- float64 %v \n", float)
	fmt.Printf("- boolean %v \n", boolean)
	fmt.Println()

	word = "hello"
	a = 200
	fmt.Println(word, a, b, c, WORLD)
}
