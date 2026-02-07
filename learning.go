//go:build learning
// +build learning

package main

import (
	"fmt"
	"sync"
)

// type Product struct {
// 	name string
// 	color Color
// 	size Size
// }

// func FilterByColor(products []Product, color Color) []*Product {
// 	result := make([]*Product, 0)

// 	for i, v := range products {
// 		if v.color = color {
// 			result = append(result, &products[i])
// 		}
// 	}
// }

func fun(s string) {
	for i := 0; i < 3; i++ {
		fmt.Println(s)

	}
}

func main() {

	ch := make(chan int, 3)
	var wg sync.WaitGroup
	// var data int

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ch <- i

			// fmt.Println(i)

		}(i)

	}

	go func() {

		wg.Wait()
		close(ch)
	}()
	for cval := range ch {
		fmt.Printf("value from ch: %v \n", cval)
	}
}

/*
func main() {
	// Your code starts here
	fmt.Println("Hello, Go!")

	// Try calling functions below
	greet("Learner")
	result := add(5, 3)
	fmt.Println("5 + 3 =", result)

	// Practice with variables
	name := "Go Developer"
	age := 25
	fmt.Printf("Name: %s, Age: %d\n", name, age)
}
*/

// // Example function: greets a person
// func greet(name string) {
// 	fmt.Printf("Hello, %s!\n", name)
// }

// // Example function: adds two numbers
// func add(a, b int) int {
// 	return a + b
// }

// // Example function: subtracts two numbers
// func subtract(a, b int) int {
// 	return a - b
// }

// // Example function: multiplies two numbers
// func multiply(a, b int) int {
// 	return a * b
// }

// Add your own functions below to practice
