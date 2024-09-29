package api

import (
	"encoding/json"
	"log"
	"net/http"

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

func HandlePunyURL(w http.ResponseWriter, r *http.Request, c *cache.Cache) {
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

	// O(n) lookup in the worst case; can this be improved?
	// Check the cache
	cachedShortID, found := c.FindByLong(req.LongURL)
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
		c.Update(existingShortID, req.LongURL)
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
	c.Update(shortID, req.LongURL)
	log.Printf("Short URL %s generated for %s", shortID, req.LongURL)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shortenResponse{ShortURL: BaseURL + shortID})
}

func HandleRedirect(w http.ResponseWriter, r *http.Request, c *cache.Cache) {
	vars := mux.Vars(r)
	shortID := vars["id"]

	// Cache lookup
	if longURL, found := c.GetLong(shortID); found {
		log.Printf("(Cache hit) %s: %s", shortID, longURL)
		http.Redirect(w, r, longURL, http.StatusFound)
		return
	}

	// If not found in cache, lookup the database
	longURL, err := db.GetLong(shortID)
	if err != nil || longURL == "" {
		log.Printf("Could not find %s in the database", shortID)
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	c.Update(shortID, longURL)
	log.Printf("(Cache miss) %s: %s has been put back into the cache.", shortID, longURL)

	http.Redirect(w, r, longURL, http.StatusFound)
}
