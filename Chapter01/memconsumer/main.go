package main

import (
        "fmt"
        "strings"
        "time"
)

func main() {
        var longStrs []string
        times := 50
        for i := 1; i <= times; i++ {
                fmt.Printf("===============%d===============\n", i)
                // each time we build a long string to consume 1MB (1000000 * 1byte) RAM
                longStrs = append(longStrs, buildString(1000000, byte(i)))
        }
        // hold the application to exit in 1 hour
        time.Sleep(3600 * time.Second)
}

// buildString build a long string with a length of `n`.
func buildString(n int, b byte) string {
        var builder strings.Builder
        builder.Grow(n)
        for i := 0; i < n; i++ {
                builder.WriteByte(b)
        }
        return builder.String()
}
