package main

import (
	"net/http"

	"log"
)

func main() {

	router := newRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
