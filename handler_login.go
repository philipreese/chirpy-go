package main

import (
	"encoding/json"
	"net/http"

	"github.com/philipreese/chirpy-go/internal/auth"
)

func (cfg *apiConfig) handlerLogin(writer http.ResponseWriter, req *http.Request) {
	type createLoginRequest struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)
	var loginRequest createLoginRequest
	if err := decoder.Decode(&loginRequest); err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Couldn't decode parameters: " + err.Error())
		return
	}

	user, err := cfg.db.GetUserByEmail(req.Context(), loginRequest.Email)
	if err != nil {
		respondWithError(writer, http.StatusUnauthorized, "Incorrect email or password")
	}

	if err := auth.CheckPasswordHash(loginRequest.Password, user.HashedPassword); err != nil {
		respondWithError(writer, http.StatusUnauthorized, "Incorrect email or password")
	}

	respondWithJSON(writer, http.StatusOK, User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	})
}