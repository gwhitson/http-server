package main

import (
    "fmt"
    "net"
)

func main() {
    conn, err := net.Dial("tcp", "localhost:8080")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer conn.Close()

    data := []byte("testtest")
    _, err = conn.Write(data)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
}
