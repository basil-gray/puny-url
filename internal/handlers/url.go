package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"puny-url/config"
	"puny-url/db"
	"puny-url/internal/cache"
	"puny-url/internal/helpers"
	"puny-url/internal/logger"

	"github.com/gorilla/mux"
)

type punyRequest struct {
	LongURL string `json:"long_url"`
}

type punyResponse struct {
	ShortURL string `json:"short_url"`
}

// Main handler functions

func HandlePunifyRequest(w http.ResponseWriter, r *http.Request, c *cache.Cache) {
	punyReq, err := parsePunyRequest(r)
	if err != nil {
		handleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortID, err := findExistingLongURL(c, punyReq.LongURL)
	if err != nil {
		handleError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if shortID != "" {
		sendPunyResponse(w, shortID)
		return
	}

	// Couldn't find ID in cache or db; generate a new one
	shortID = helpers.GenerateShortID()

	// Save to cache and db
	c.Update(shortID, punyReq.LongURL)
	db.StoreURL(shortID, punyReq.LongURL)
	logger.Info(fmt.Sprintf("Short URL %s generated for %s", shortID, punyReq.LongURL))
	sendPunyResponse(w, shortID)
}

func HandleRedirect(w http.ResponseWriter, r *http.Request, c *cache.Cache) {
	vars := mux.Vars(r)
	shortID := vars["id"]

	if !helpers.IsValidShortID(shortID) {
		handleError(w, "Invalid ID", http.StatusBadRequest)
	}

	// Cache lookup
	if longURL, found := c.Load(shortID); found {
		c.Update(shortID, longURL)
		logger.Info(fmt.Sprintf("(Cache hit) Redirecting to %s: %s", shortID, longURL))
		http.Redirect(w, r, longURL, http.StatusFound)
		return
	}

	// If not found in cache, lookup the database
	longURL, err := db.GetLong(shortID)
	if err != nil || longURL == "" {
		handleError(w, "URL not found", http.StatusNotFound)
		return
	}

	c.Update(shortID, longURL)
	logger.Info(fmt.Sprintf("(Cache miss) Redirecting to %s: %s. URL has been put back into the cache.", shortID, longURL))

	http.Redirect(w, r, longURL, http.StatusFound)
}

// Helper functions

func parsePunyRequest(r *http.Request) (*punyRequest, error) {
	if r.ContentLength == 0 {
		return nil, fmt.Errorf("empty request")
	}

	var req punyRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, fmt.Errorf("invalid request: %v", err)
	}

	if req.LongURL == "" {
		return nil, fmt.Errorf("long URL is empty")
	}

	req.LongURL, err = helpers.ParseURL(req.LongURL)
	if err != nil {
		return nil, err
	}

	return &req, nil
}

func sendPunyResponse(w http.ResponseWriter, shortId string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(punyResponse{ShortURL: config.BaseURL() + shortId})
}

func handleError(w http.ResponseWriter, errMsg string, statusCode int) {
	logger.Error(errMsg)
	http.Error(w, errMsg, statusCode)
}

func findExistingLongURL(c *cache.Cache, longURL string) (string, error) {
	// Check cache
	cachedShortID, found := c.FindByLong(longURL)
	if found {
		c.Update(cachedShortID, longURL)
		logger.Info(fmt.Sprintf("(Cache hit) Shortern URL attempt received but long URL was found in cache: %s: %s", longURL, cachedShortID))
		return cachedShortID, nil
	}

	// If not in the cache, check the database
	existingShortID, err := db.GetShortIDByLongURL(longURL)
	if err != nil {
		return "", fmt.Errorf("error checking for existing URL: %w", err)
	}

	if existingShortID != "" {
		c.Update(existingShortID, longURL)
		logger.Info(fmt.Sprintf("(DB) %s: %s Existing short URL retrieved from the database and has been put back into cache.", longURL, existingShortID))
		return existingShortID, nil
	}

	return "", nil
}
