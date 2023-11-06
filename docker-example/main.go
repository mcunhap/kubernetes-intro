package main

import (
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "Hello World")
}

func main() {
	http.HandleFunc("/hello", helloHandler)

	fmt.Println("Server is running on http://localhost:8080/hello")
	http.ListenAndServe(":8080", nil)
}
