package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name          string
		password      string
		hash          string
		wantErr       bool
		matchPassword bool
	}{
		{
			name:          "Correct password",
			password:      password1,
			hash:          hash1,
			wantErr:       false,
			matchPassword: true,
		},
		{
			name:          "Incorrect password",
			password:      "wrongPassword",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Password doesn't match different hash",
			password:      password1,
			hash:          hash2,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Empty password",
			password:      "",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Invalid hash",
			password:      password1,
			hash:          "invalidhash",
			wantErr:       true,
			matchPassword: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && match != tt.matchPassword {
				t.Errorf("CheckPasswordHash() expects %v, got %v", tt.matchPassword, match)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "my-secret"

	validToken, _ := MakeJWT(userID, secret, time.Hour)
	expiredToken, _ := MakeJWT(userID, secret, -time.Hour)

	tests := []struct {
		name       string
		token      string
		secret     string
		wantErr    bool
		expectedID uuid.UUID
	}{
		{
			name:       "Valid token",
			token:      validToken,
			secret:     secret,
			wantErr:    false,
			expectedID: userID,
		},
		{
			name:       "Expired token",
			token:      expiredToken,
			secret:     secret,
			wantErr:    true,
			expectedID: uuid.Nil,
		},
		{
			name:       "Invalid secret",
			token:      validToken,
			secret:     "wrong-secret",
			wantErr:    true,
			expectedID: uuid.Nil,
		},
		{
			name:       "Malformed token",
			token:      "not-a-valid-jwt",
			secret:     secret,
			wantErr:    true,
			expectedID: uuid.Nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, err := ValidateJWT(tt.token, tt.secret)

			if (err != nil) != tt.wantErr {
				t.Errorf(
					"ValidateJWT() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}

			if !tt.wantErr && gotID != tt.expectedID {
				t.Errorf(
					"ValidateJWT() expected ID = %v, got %v",
					tt.expectedID,
					gotID,
				)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name       string
		headers    http.Header
		wantErr    bool
		wantToken  string
	}{
		{
			name: "Valid bearer token",
			headers: http.Header{
				"Authorization": []string{"Bearer abc123"},
			},
			wantErr:   false,
			wantToken: "abc123",
		},
		{
			name: "Missing authorization header",
			headers: http.Header{},
			wantErr: true,
		},
		{
			name: "Invalid authorization format",
			headers: http.Header{
				"Authorization": []string{"Basic abc123"},
			},
			wantErr: true,
		},
		{
			name: "Authorization header without token",
			headers: http.Header{
				"Authorization": []string{"Bearer "},
			},
			wantErr:   false,
			wantToken: "",
		},
		{
			name: "Invalid bearer prefix",
			headers: http.Header{
				"Authorization": []string{"bearer abc123"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, err := GetBearerToken(tt.headers)

			if (err != nil) != tt.wantErr {
				t.Errorf(
					"GetBearerToken() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}

			if !tt.wantErr && gotToken != tt.wantToken {
				t.Errorf(
					"GetBearerToken() token = %v, wantToken %v",
					gotToken,
					tt.wantToken,
				)
			}
		})
	}
}