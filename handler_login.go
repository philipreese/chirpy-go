package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/philipreese/chirpy-go/internal/auth"
	"github.com/philipreese/chirpy-go/internal/database"
)

func (cfg *apiConfig) handlerLogin(writer http.ResponseWriter, req *http.Request) {
	type createLoginRequest struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
	}

	type loginResponse struct {
		User		
	    Token        string `json:"token"`
		RefreshToken string `json:"refresh_token,omitempty"`
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
		return
	}

	if err := auth.CheckPasswordHash(loginRequest.Password, user.HashedPassword); err != nil {
		respondWithError(writer, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	tokenString, err := auth.MakeJWT(user.ID, cfg.tokenSecret, time.Hour)
	if err != nil {		
		respondWithError(writer, http.StatusInternalServerError, "Failed to create token: " + err.Error())
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Failed to create refresh token: " + err.Error())
		return
	}

	_, err = cfg.db.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token: refreshToken,
		UserID: user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
		RevokedAt: sql.NullTime{},
	})
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Failed to save refresh token: " + err.Error())
		return
	}

	respondWithJSON(writer, http.StatusOK, loginResponse{
		User: User{
			ID: user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email: user.Email,
		},
		Token: tokenString,
		RefreshToken: refreshToken,
	})
}