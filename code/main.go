package main

import (
    "fmt"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Hello, World!")
}

func main() {
    fmt.Println("Server Start")
    http.HandleFunc("/", handler)
    http.ListenAndServe(":80", nil)
}

