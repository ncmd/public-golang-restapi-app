package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {

	// hello world
	http.HandleFunc("/", helloworld)

	// response pong
	http.HandleFunc("/api/ping", pingpong)

	// listen on socket
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // Setting a Default port to 8000 to be used locally
	}
	fmt.Println(http.ListenAndServe(":"+port, nil))
	http.ListenAndServe(":"+port, nil)

}
