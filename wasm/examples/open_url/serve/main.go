package main

import (
	"log"
	"net/http"
)

func main() {
	port := "8080"
	log.Println("listen on http://localhost:", port)
	dir := "./dev.flowingspdg.wasm.sdPlugin/pi"
	log.Println("Serving on :", dir)
	http.Handle("/", http.FileServer(http.Dir(dir)))
	http.ListenAndServe(":"+port, nil)
}
