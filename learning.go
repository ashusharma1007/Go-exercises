//go:build learning
// +build learning

package main

import "fmt"

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

// func ch1e(ch chan string, s string) {
// 	ch <- s
// }

// func ch2o(ch chan string, s string) {
// 	ch <- s
// }

func main() {

	ch1 := make(chan string, 10)
	ch2 := make(chan string, 10)

	go func() {
		for i := 0; i < 5; i++ {
			ch1 <- "ping"
			ch1 <- "pong"

		}
		close(ch1)

	}()

	go func() {
		for val := range ch1 {
			ch2 <- val

		}
		close(ch2)
	}()

	for val := range ch2 {
		fmt.Println(val)
	}
}
