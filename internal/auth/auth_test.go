package auth

import (
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