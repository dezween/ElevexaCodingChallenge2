package main

import (
	"log"
	"net/http"

	"github.com/dezween/ElevexaCodingChallenge2/internal/config"
	"github.com/dezween/ElevexaCodingChallenge2/internal/server"
)

func main() {
	cfg := config.LoadConfig()
	router := server.NewRouter()
	log.Printf("Kyber Transit API server running on %s", cfg.Port)
	log.Fatal(http.ListenAndServe(cfg.Port, router))
}
