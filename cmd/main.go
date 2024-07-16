package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/hello-world", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world!"))
	})

	fmt.Println("Server listening at port :8080 🚀")
	http.ListenAndServe(":8080", nil)
}
