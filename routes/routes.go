package routes

import (
	"puny-url/api"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/shorten", api.PunifyURLHandler).Methods("POST")
	r.HandleFunc("/{id}", api.RedirectHandler).Methods("GET")

	return r
}
