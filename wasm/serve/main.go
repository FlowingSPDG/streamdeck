package main

import (
	"log"
	"net/http"
)

func main() {
	port := "8080"
	log.Printf("listen on http://localhost:%s", port)
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.ListenAndServe(":"+port, nil)
}
