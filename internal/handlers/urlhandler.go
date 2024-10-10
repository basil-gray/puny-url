package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"puny-url/config"
	"puny-url/internal/cache"
	"puny-url/internal/helpers"
	"puny-url/internal/logger"
	"puny-url/internal/models"
	"puny-url/internal/services"

	"github.com/gorilla/mux"
)

func HandlePunifyRequest(w http.ResponseWriter, r *http.Request, c *cache.Cache) {
	punyReq, err := services.ParsePunyRequest(r)
	if err != nil {
		handleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortID, err := services.GetShortID(c, punyReq.LongURL)
	if err != nil {
		handleError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendPunyResponse(w, shortID)
}

func HandleRedirect(w http.ResponseWriter, r *http.Request, c *cache.Cache) {
	vars := mux.Vars(r)
	shortID := vars["id"]

	if !helpers.IsValidShortID(shortID) {
		handleError(w, "Invalid ID", http.StatusBadRequest)
	}

	longURL, err := services.GetLongId(c, shortID)
	if err != nil {
		handleError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info(fmt.Sprintf("Redirecting to %s: %s", shortID, longURL))

	http.Redirect(w, r, longURL, http.StatusFound)
}

func sendPunyResponse(w http.ResponseWriter, shortId string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.PunyResponse{ShortURL: config.BaseURL() + shortId})
}

func handleError(w http.ResponseWriter, errMsg string, statusCode int) {
	logger.Error(errMsg)
	http.Error(w, errMsg, statusCode)
}
