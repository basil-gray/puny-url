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

func HandlePunyURL(w http.ResponseWriter, r *http.Request) {
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

	now := time.Now()

	// O(n) lookup in the worst case; can this be improved?
	// Check the cache
	cachedShortID, found := cache.Find(req.LongURL, now)
	if found {
		log.Printf("(Cache hit) Short URL found for %s: %s", req.LongURL, cachedShortID)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(shortenResponse{ShortURL: BaseURL + cachedShortID})
		return
	}

	// Check the database
	existingShortID, err := db.GetShortIDByLongURL(req.LongURL)
	if err != nil {
		log.Printf("Error checking for existing URL: %s", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if existingShortID != "" {
		log.Printf("Existing short URL found in the database for %s: %s", req.LongURL, existingShortID)
		cache.Update(existingShortID, req.LongURL, now)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(shortenResponse{ShortURL: BaseURL + existingShortID})
		return
	}

	// Generate new short ID and store in database
	shortID := helpers.GenerateShortID()
	if err := db.StoreURL(shortID, req.LongURL); err != nil {
		log.Printf("An error has occurred when storing the URL: %s", err.Error())
		http.Error(w, "Failed to store URL", http.StatusInternalServerError)
		return
	}

	// Update the cache
	cache.Update(shortID, req.LongURL, now)
	log.Printf("Short URL %s generated for %s", shortID, req.LongURL)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shortenResponse{ShortURL: BaseURL + shortID})
}

func HandleRedirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortID := vars["id"]

	now := time.Now()

	if entry, found := cache.Map.Load(shortID); found {
		longURL := entry.(cache.CacheEntry).LongURL

		// Update the last accessed time
		log.Printf("(Cache hit) %s: %s", shortID, longURL)
		cache.Map.Store(shortID, cache.CacheEntry{LongURL: longURL, Timestamp: entry.(cache.CacheEntry).Timestamp, LastAccessed: now})
		http.Redirect(w, r, longURL, http.StatusFound)
		return
	}

	// If not found in cache, lookup the database
	longURL, err := db.GetLongURL(shortID)
	if err != nil || longURL == "" {
		log.Printf("Could not find %s in the database", shortID)
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	cache.Update(shortID, longURL, now)
	log.Printf("(Cache miss) %s: %s has been put back into the cache.", shortID, longURL)

	http.Redirect(w, r, longURL, http.StatusFound)
}
