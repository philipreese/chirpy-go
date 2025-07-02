package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct password",
			password: password1,
			hash:     hash1,
			wantErr:  false,
		},
		{
			name:     "Incorrect password",
			password: "wrongPassword",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Password doesn't match different hash",
			password: password1,
			hash:     hash2,
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Invalid hash",
			password: password1,
			hash:     "invalidhash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "testsecret"
	validToken, _ := MakeJWT(userID, tokenSecret, time.Hour)

	tests := []struct {
		name           string
		tokenString    string
		tokenSecret    string
		expectedUserID uuid.UUID
		expectedErr    bool
	}{
		{
			name: "Valid token",
			tokenString: validToken,
			tokenSecret: tokenSecret,
			expectedUserID: userID,
			expectedErr: false,
		},
		{
			name: "Invalid token",
			tokenString: "test-invalid.token",
			tokenSecret: tokenSecret,
			expectedUserID: uuid.Nil,
			expectedErr: true,
		},
		{
			name: "Invalid secret",
			tokenString: validToken,
			tokenSecret: "something-else",
			expectedUserID: uuid.Nil,
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsedID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.expectedErr {
				t.Errorf("failed to validate JWT: %v", err)
				return
			}

			if parsedID != tt.expectedUserID {
				t.Errorf("expected userID %v, got %v", tt.expectedUserID, parsedID)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name string
		headerString string
		expectedErr bool
	}{
		{
			name: "Valid bearer token",
			headerString: "Bearer abcdef",
			expectedErr: false,
		},
		{
			name: "Missing authorization header",
			headerString: "",
			expectedErr: true,
		},
		{
			name: "Invalid authorization header format",
			headerString: "Invalid abcdef",
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := http.Header{}
			if tt.headerString != "" {
				headers.Set("Authorization", tt.headerString)
			}
			_, err := GetBearerToken(headers)
			if (err != nil) != tt.expectedErr {
				t.Errorf("GetBearerToken() error = %v, expectedErr %v", err, tt.expectedErr)
			}
		})
	}
}