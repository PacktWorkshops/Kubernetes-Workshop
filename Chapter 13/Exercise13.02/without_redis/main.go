package main

import (
	"fmt"
	"log"
	"net/http"
)

var num = 0
func main() {
	fmt.Println("Starting HTTP Server")
	http.HandleFunc("/get-number", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET"{
			num++
			fmt.Fprintf(w, "{number: %d}", num)
		} else {
			w.WriteHeader(400)
			fmt.Fprint(w, "{\"error\": \"Only GET HTTP method is supported.\"}")
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
