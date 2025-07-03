package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/philipreese/chirpy-go/internal/auth"
	"github.com/philipreese/chirpy-go/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
}

type userRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) handlerCreateUser(writer http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var userRequest userRequest
	if err := decoder.Decode(&userRequest); err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Couldn't decode parameters: " + err.Error())
		return
	}

	hashedPassword, err := auth.HashPassword(userRequest.Password)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Couldn't hash password: " + err.Error())
		return
	}

	dbUser, err := cfg.db.CreateUser(req.Context(), database.CreateUserParams{
		Email: userRequest.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Couldn't create user: " + err.Error())
		return
	}
	
	user := User{
		ID: dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email: dbUser.Email,
	}
	respondWithJSON(writer, http.StatusCreated, user)
}

func (cfg *apiConfig) handlerUpdateUser(writer http.ResponseWriter, req *http.Request) {
	tokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(writer, http.StatusUnauthorized, "Couldn't get bearer token: " + err.Error())
		return
	}

	userID, err := auth.ValidateJWT(tokenString, cfg.tokenSecret)
	if err != nil {
		respondWithError(writer, http.StatusUnauthorized, "Couldn't validate JWT: " + err.Error())
		return
	}

	decoder := json.NewDecoder(req.Body)
	var userRequest userRequest
	if err := decoder.Decode(&userRequest); err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Couldn't decode parameters: " + err.Error())
		return
	}

	hashedPassword, err := auth.HashPassword(userRequest.Password)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Couldn't hash password: " + err.Error())
		return
	}

	user, err := cfg.db.UpdateUser(req.Context(), database.UpdateUserParams{
		ID: userID,
		Email: userRequest.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Failed to update user: " + err.Error())
		return
	}

	respondWithJSON(writer, http.StatusOK, User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	})
}