package main

import "fmt"
import "rsc.io/quote"

func main() {
	fmt.Println("hello go start")
	fmt.Println(quote.Hello())
	fmt.Println(quote.Go())
	fmt.Println("hello go end")
}
