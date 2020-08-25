package main

import (
        "fmt"
        "log"
        "net/http"

        "github.com/go-redis/redis/v7"
)

var dbClient *redis.Client
var key = "pv"

func init() {
        dbClient = redis.NewClient(&redis.Options{
                Addr:     "db:6379",
        })
}

func main() {
        http.HandleFunc("/", handler)
        log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
        log.Printf("Ping from %s", r.RemoteAddr)
        pageView, err := dbClient.Incr(key).Result()
        if err != nil {
                panic(err)
        }
        fmt.Fprintf(w, "Hello, you're visitor #%v.\n", pageView)
}
