package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"puny-url/cache"
	"puny-url/db"
	"puny-url/helpers"

	"github.com/gorilla/mux"
)

type shortenRequest struct {
	LongURL string `json:"long_url"`
}

type shortenResponse struct {
	ShortURL string `json:"short_url"`
}

const (
	BaseURL = "http://localhost:8080/"
)

func PunifyURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength == 0 {
		http.Error(w, "Empty request", http.StatusBadRequest)
		return
	}
	var req shortenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("An error has occurred when decoding the URL: %s", err.Error())
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.LongURL == "" {
		log.Printf("An error has occurred when decoding the UTL: long URL is empty")
		http.Error(w, "Long URL is empty", http.StatusBadRequest)
		return
	}

	req.LongURL, err = helpers.ValidateURL(req.LongURL)

	if err != nil {
		log.Printf("An error has occurred when validating the URL: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortID := helpers.GenerateShortID()

	// Store the URL in the database
	if err := db.StoreURL(shortID, req.LongURL); err != nil {
		log.Printf("An error has occurred when storing the URL: %s", err.Error())
		http.Error(w, "Failed to store URL", http.StatusInternalServerError)
		return
	}

	// Update the cache
	now := time.Now()
	cache.Map.Store(shortID, cache.CacheEntry{Value: req.LongURL, Timestamp: now, LastAccessed: now})
	log.Printf("Short URL %s generated for %s", shortID, req.LongURL)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shortenResponse{ShortURL: BaseURL + shortID})
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortID := vars["id"]

	now := time.Now()

	if entry, found := cache.Map.Load(shortID); found {
		longURL := entry.(cache.CacheEntry).Value

		// Update the last accessed time
		log.Printf("(Cache hit) %s: %s", shortID, longURL)
		cache.Map.Store(shortID, cache.CacheEntry{Value: longURL, Timestamp: entry.(cache.CacheEntry).Timestamp, LastAccessed: now})
		http.Redirect(w, r, longURL, http.StatusFound)
		return
	}

	// If not found in cache, lookup the database
	longURL, err := db.GetLongURL(shortID)
	if err != nil || longURL == "" {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	cache.Map.Store(shortID, cache.CacheEntry{Value: longURL, Timestamp: now, LastAccessed: now})
	log.Printf("(Cache miss) %s: %s has been put back into the cache.", shortID, longURL)

	http.Redirect(w, r, longURL, http.StatusFound)
}
