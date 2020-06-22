package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
)

var (
	num int
)

const (
	t1      = "UPDATE counter SET num = @num := num + 1 WHERE id=1;"
	t2      = "SELECT num from counter where id=?"
	startup = "CREATE TABLE IF NOT EXISTS counter(id INT, num INT);"
	startupEntry = "INSERT IGNORE INTO db1.counter VALUES (1,0);"
)

func main() {
	fmt.Println("Starting MySQL Connection")
	db, err := sql.Open("mysql", os.Getenv("CONN_STRING"))
	defer db.Close()
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(startup)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(startupEntry)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting HTTP server")
	http.HandleFunc("/get-number", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			FibonacciLoop(10000000)
			tx, err := db.Begin()
			if err != nil {
				panic(err)
			}
			_, err = tx.Exec(t1)
			if err != nil {
				tx.Rollback()
				fmt.Println(err)
			}
			err = tx.Commit()
			if err != nil {
				fmt.Println(err)
			}
			row := db.QueryRow(t2, 1)
			switch err := row.Scan(&num); err {
			case sql.ErrNoRows:
				fmt.Println("No rows were returned!")
			case nil:
				fmt.Fprintf(w, "{number: %d}\n", num)
			default:
				panic(err)
			}
		} else {
			w.WriteHeader(400)
			fmt.Fprint(w, "{\"error\": \"Only GET HTTP method is supported.\"}")
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

//Taken from: https://www.thepolyglotdeveloper.com/2016/12/fibonacci-sequence-printed-golang/
func FibonacciLoop(n int) int {
	f := make([]int, n+1, n+2)
	if n < 2 {
		f = f[0:2]
	}
	f[0] = 0
	f[1] = 1
	for i := 2; i <= n; i++ {
		f[i] = f[i-1] + f[i-2]
	}
	return f[n]
}