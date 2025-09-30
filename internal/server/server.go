package server

import (
	"github.com/dezween/ElevexaCodingChallenge2/internal/handlers"
	"github.com/dezween/ElevexaCodingChallenge2/internal/routes"
	"github.com/gorilla/mux"
)

// NewRouter returns a fully configured HTTP router for the Kyber Transit API.
// Routes are named to allow URL building via mux.Route.URL in tests and other code.
func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc(routes.RouteCreateKey, handlers.CreateKeyHandler).Methods("POST").Name(routes.RouteNameCreateKey)
	r.HandleFunc(routes.RouteEncrypt, handlers.EncryptHandler).Methods("POST").Name(routes.RouteNameEncrypt)
	r.HandleFunc(routes.RouteDecrypt, handlers.DecryptHandler).Methods("POST").Name(routes.RouteNameDecrypt)
	r.HandleFunc("/health", handlers.HealthHandler).Methods("GET")
	return r
}
