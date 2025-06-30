package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/philipreese/chirpy-go/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(writer http.ResponseWriter, req *http.Request) {
	type chirpRequest struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(req.Body)
	var chirpReq chirpRequest
	if err := decoder.Decode(&chirpReq); err != nil {
		respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	if len(chirpReq.Body) > 400 {
		respondWithError(writer, http.StatusBadRequest, "Chirp is too long")
		return
	}
	cleanedBody := getCleanedBody(chirpReq.Body)

	dbChirp, err := cfg.db.CreateChirp(req.Context(), database.CreateChirpParams{
		Body: cleanedBody,
		UserID: chirpReq.UserID,
	})
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(writer, http.StatusCreated, Chirp{
		ID: dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body: dbChirp.Body,
		UserID: dbChirp.UserID,
	})
}

func (cfg *apiConfig) handlerGetChirps(writer http.ResponseWriter, req *http.Request) {
	dbChirps, err := cfg.db.GetChirps(req.Context())
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID: dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body: dbChirp.Body,
			UserID: dbChirp.UserID,
		})
	}

	respondWithJSON(writer, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerGetChirpByID(writer http.ResponseWriter, req *http.Request) {
	chirpID, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	dbChirp, err := cfg.db.GetChirpByID(req.Context(), chirpID)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(writer, http.StatusOK, Chirp{
		ID: dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body: dbChirp.Body,
		UserID: dbChirp.UserID,
	})
}

func getCleanedBody(text string) string {
	profanityList := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(text, " ")
	for i, word := range(words) {
		for _, profanity := range profanityList {
			if strings.ToLower(word) == profanity {
				words[i] = "****"
			}
		}
	}
	return strings.Join(words, " ")
}