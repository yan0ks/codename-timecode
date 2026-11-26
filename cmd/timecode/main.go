package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", health)

	addr := "localhost:8080"
	log.Printf("Timecode: http://%s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func health(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("ok\n"))
}
