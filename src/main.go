package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

func main() {
	fmt.Printf("Hello, world")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "JGP.SH")
	})

	log.Println("Creating tcp listener on port 8080")
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Serving server on listener")
	http.Serve(listener, nil)
}
