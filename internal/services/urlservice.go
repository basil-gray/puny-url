package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"puny-url/db"
	"puny-url/internal/cache"
	"puny-url/internal/helpers"
	"puny-url/internal/logger"
	"puny-url/internal/models"
)

func ParsePunyRequest(r *http.Request) (*models.PunyRequest, error) {
	if r.ContentLength == 0 {
		return nil, fmt.Errorf("empty request")
	}

	var req models.PunyRequest
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

func GetShortID(c *cache.Cache, longURL string) (string, error) {
	shortId, found := c.FindByLong(longURL)
	if found {
		logger.Info(fmt.Sprintf("(Cache hit) %s: %s", shortId, longURL))
		c.UpdateCache(shortId, longURL)
		return shortId, nil
	}

	// If not in the cache, check the database
	shortId, err := db.GetShortIDByLongURL(longURL)
	if err != nil {
		errorMsg := fmt.Errorf("error checking for existing URL: %w", err)
		logger.Info(errorMsg.Error())
		return "", errorMsg
	}

	if shortId != "" {
		c.UpdateCache(shortId, longURL)
		logger.Info(fmt.Sprintf("(Cache update) %s: %s retrieved from DB", shortId, longURL))
		return shortId, nil
	}

	shortId = helpers.GenerateShortID()
	err = db.CreateURL(shortId, longURL)
	logger.Info(fmt.Sprintf("New short URL generated. %s: %s", shortId, longURL))
	if err != nil {
		return "", fmt.Errorf("failed to save new URL %s", longURL)
	}

	c.UpdateCache(shortId, longURL)

	return shortId, nil
}

func GetLongId(c *cache.Cache, shortID string) (string, error) {
	if longURL, found := c.LoadFromCache(shortID); found {
		logger.Info(fmt.Sprintf("(Cache hit) %s, %s", shortID, longURL))
		c.UpdateCache(shortID, longURL)
		return longURL, nil
	}

	logger.Info(fmt.Sprintf("(Cache miss) %s", shortID))

	longURL, err := db.LoadFromDB(shortID)
	if err != nil {
		return "", err
	}

	if longURL == "" {
		return "", fmt.Errorf("could not find long URL for %s", shortID)
	}

	c.UpdateCache(shortID, longURL)

	return longURL, nil
}
