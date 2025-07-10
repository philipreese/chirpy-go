package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/philipreese/chirpy-go/internal/auth"
)

type webhookEvent struct {
	Event string `json:"event"`
	Data struct {
		UserID uuid.UUID `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) handlerWebhook(writer http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var webhookEvent webhookEvent
	if err := decoder.Decode(&webhookEvent); err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Couldn't decode parameters: " + err.Error())
		return
	}

	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		respondWithError(writer, http.StatusUnauthorized, "Couldn't get API key: " + err.Error())
		return
	}

	if cfg.polkaKey != apiKey {
		respondWithError(writer, http.StatusUnauthorized, "Invalid API key")
		return
	}

	if webhookEvent.Event != "user.upgraded" {
		writer.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.db.UpgradeUser(req.Context(), webhookEvent.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(writer, http.StatusNotFound, "User not found: " + err.Error())
			return
		}
		respondWithError(writer, http.StatusNotFound, "Couldn't update user: " + err.Error())
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}