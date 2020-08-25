package main

import (
        "fmt"
        "log"
        "net/http"
)

var pageView int64

func main() {
        http.HandleFunc("/", handler)
        log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
        log.Printf("Ping from %s", r.RemoteAddr)
        pageView++
        fmt.Fprintf(w, "Hello, you're visitor #%d !", pageView)
}
