package main

import (
	"log"
	"net/http"
)

func main() {
	port := "8080"
	log.Printf("listen on http://localhost:%s\n", port)
	dir := "./dev.flowingspdg.binding.sdPlugin/pi"
	log.Println("Serving on :", dir)
	http.Handle("/", http.FileServer(http.Dir(dir)))
	http.ListenAndServe(":"+port, nil)
}
