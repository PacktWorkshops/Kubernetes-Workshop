package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"net/http"
	"strconv"
)

var (
	rdb *redis.Client
)

func main() {
	fmt.Println("Starting Redis Connection")
	client := redis.NewClient(&redis.Options{
		Addr:     "redis.default:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	fmt.Println("Starting HTTP server")
	http.HandleFunc("/get-number", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			val, err := client.Get("num").Result()
			if err == redis.Nil {
				fmt.Println("num does not exist")
				err := client.Set("num", "0", 0).Err()
				if err != nil {
					panic(err)
				}
			} else if err != nil {
				w.WriteHeader(500)
				panic(err)
			} else {
				fmt.Println("num", val)
				num, err := strconv.Atoi(val)
				if err != nil {
					w.WriteHeader(500)
					fmt.Println(err)
				} else {
					num++
					err := client.Set("num", strconv.Itoa(num), 0).Err()
					if err != nil {
						panic(err)
					}
					fmt.Fprintf(w, "{number: %d}", num)
				}
			}
		} else {
			w.WriteHeader(400)
			fmt.Fprint(w, "{\"error\": \"Only GET HTTP method is supported.\"}")
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

