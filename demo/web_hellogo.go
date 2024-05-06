package main

import (
	"fmt"
	_ "fmt"
	"net/http"
	_ "net/http"
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, "web hello go")
	})
	err := http.ListenAndServe(":8880", nil)
	if err != nil {
		return
	}

}
