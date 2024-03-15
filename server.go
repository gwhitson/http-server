package main

import (
	//"unicode/utf8"
	"bytes"
	"fmt"
	"net"
    "os"
)

func main() {
    listener, err := net.Listen("tcp","localhost:8080")
    check(err)
    defer listener.Close()

    fmt.Println("Listening on port 8080")

    for {
        conn, err := listener.Accept()
        check(err)
        defer conn.Close()
        
        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
    buffer := make([]byte, 1024)
    length, _ := conn.Read(buffer)

    t, r, _ := getRequest(buffer, length)
    t, r = bytes.Trim(t, "\x00"), bytes.Trim(r, "\x00")


    rtype := string(t)
    res := string(r)

    if res[len(res) - 1] == '/' {
        res += "index.html"
    }
    if res[0] == '/' {
        res = "." + res
    }

    switch rtype {
    case "GET":
        fmt.Println("serving request")
        dat, err := os.ReadFile(res)
        check(err)
        res := "HTTP/1.1 200 OK\r\nContent-Type: text/html\r\nhtdocs\r\n"
        res += string(dat)
        conn.Write([]byte(res))
    default:
        conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
    }

    defer conn.Close()
}

func printBytes(buffer []byte) {
    fmt.Printf("%s", buffer[:])
}

func getRequest(buffer []byte, length int) ([]byte, []byte, []byte) {
    rtype := make([]byte, 16)
    res := make([]byte, 32)
    rhttp := make([]byte, 16)
    ind := 0
    i := buffer[ind]

    // get type of request
    for i != byte(' ') {
        rtype = append(rtype[:], i)
        ind = ind + 1
        i = buffer[ind]
    }

    // get resource requested
    ind = ind + 1
    i = buffer[ind]
    for i != byte(' ') {
        res = append(res[:], i)
        ind = ind + 1
        i = buffer[ind]
    }

    // get HTTP protocol
    ind = ind + 1
    i = buffer[ind]
    for i != byte('\n') {
        rhttp = append(rhttp[:], i)
        ind = ind + 1
        i = buffer[ind]
    }
    return rtype, res, rhttp
}

func check(e error) {
    if e != nil {
        fmt.Println(e)
    }
}
