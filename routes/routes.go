package routes

import (
	"net/http"
	"puny-url/internal/cache"
	"puny-url/internal/handlers"

	"github.com/gorilla/mux"
)

func withCache(c *cache.Cache, handler func(http.ResponseWriter, *http.Request, *cache.Cache)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, c)
	}
}

func SetupRoutes(c *cache.Cache) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/shorten", withCache(c, handlers.HandlePunifyRequest)).Methods("POST")
	r.HandleFunc("/{id}", withCache(c, handlers.HandleRedirect)).Methods("GET")

	return r
}
