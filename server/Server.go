package main

import (
    "net/http"
    "log"
		//"fmt"
)

func main() {
	router := NewRouter()
    log.Fatal(http.ListenAndServe(":8080", router))
}

