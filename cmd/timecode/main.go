package main

import (
	"log"
	"net/http"

	"github.com/yan0ks/codename-timecode/internal/auth"
)

func main() {
	mux := http.NewServeMux()

	repo := auth.NewMemoryRepo()

	authService := auth.NewService(repo)

	authHandler := auth.NewHandler(authService)

	mux.HandleFunc("POST /login", authHandler.Login)
	mux.HandleFunc("POST /logout", authHandler.Logout)
	mux.HandleFunc("GET /health", health)

	mux.Handle("GET /me", authService.Middleware(http.HandlerFunc(authHandler.Me)))

	addr := "localhost:8080"
	log.Printf("Timecode: http://%s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func health(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("ok\n"))
}
