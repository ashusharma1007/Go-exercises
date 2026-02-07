//go:build learning
// +build learning

package main

import (
	"fmt"
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

// func fun(s string) {
// 	for i := 0; i < 3; i++ {
// 		fmt.Println(s)

// 	}
// }

// func genMSg(ch <-chan string){
// 	msg := "get the msg"

// }

func ch1e(ch chan int, n int) {
	if n%2 == 0 {
		ch <- n
	}
}

func ch2o(ch chan int, n int) {
	if n%2 == 1 {
		ch <- n
	}
}

func main() {

	ch1 := make(chan int, 5)
	ch2 := make(chan int, 5)

	go func() {
		for i := 1; i < 11; i++ {
			ch1e(ch1, i)
			ch2o(ch2, i)
		}
		close(ch1)
		close(ch2)
	}()

	for i := 0; i < 10; i++ {
		select {
		case m1, ok := <-ch1:
			if ok {
				fmt.Println("Even:", m1)
			}
		case m2, ok := <-ch2:
			if ok {
				fmt.Println("Odd:", m2)
			}
		}
	}
}
