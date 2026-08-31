package auth

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	// TokenTypeAccess -
	TokenTypeAccess TokenType = "chirpy-access"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Issuer:    string(TokenTypeAccess),
		Subject:   userID.String(),
	})
	token,err := t.SignedString([]byte(tokenSecret))
	if err!=nil{
		log.Printf("Error occured while signing token: %v", err)
		return "", err
	}
	return token ,nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error){
	claims:= &jwt.RegisteredClaims{}
	token, err:= jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{},error){
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(tokenSecret), nil
		},
	)
	if err != nil {
		return uuid.Nil, err
	}
	if !token.Valid{
		return uuid.Nil, fmt.Errorf("invalid token")
	}
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, err
	}
	issuer, err := claims.GetIssuer()
	if issuer != string(TokenTypeAccess) {
		return uuid.Nil, fmt.Errorf("invalid issuer")
	}

	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil

}