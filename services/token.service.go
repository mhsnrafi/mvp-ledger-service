package services

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	db "ledger-service/models"
	"math/rand"
	"strconv"
	"time"
)

// CreateToken create a new token record
func CreateToken(email string, tokenType string, expiresAt time.Time) (db.Token, error) {
	// Generate a random UUID
	rand.Seed(time.Now().UnixNano())
	ID := rand.Int63()
	claims := &db.UserClaims{
		Email: email,
		Type:  tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Subject:   strconv.FormatInt(ID, 10),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(Config.JWTSecretKey))
	if err != nil {
		return db.Token{}, errors.New("cannot create access token")
	}

	tokenModel := db.Token{
		ID:          ID,
		Token:       tokenString,
		Type:        tokenType,
		ExpiresAt:   expiresAt,
		Blacklisted: false,
	}

	if err := DbConnection.Create(&tokenModel).Error; err != nil {
		return db.Token{}, fmt.Errorf("cannot save access token to db %v", err)
	}

	return tokenModel, nil
}

// DeleteTokenById delete token with id
func DeleteTokenById(tokenId int64) error {
	token := &db.Token{}

	if err := DbConnection.Where("id = ?", tokenId).First(&token).Error; err != nil {
		return err
	}
	return DbConnection.Delete(token).Error
}

// GenerateAccessTokens generates "access" and "refresh" token for user
func GenerateAccessTokens(email string) (db.Token, db.Token, error) {
	accessExpiresAt := time.Now().Add(time.Duration(Config.JWTAccessExpirationMinutes) * time.Minute)
	refreshExpiresAt := time.Now().Add(time.Duration(Config.JWTRefreshExpirationDays) * time.Hour * 24)

	accessToken, err := CreateToken(email, db.TokenTypeAccess, accessExpiresAt)
	if err != nil {
		return db.Token{}, db.Token{}, err
	}

	refreshToken, err := CreateToken(email, db.TokenTypeRefresh, refreshExpiresAt)
	if err != nil {
		return db.Token{}, db.Token{}, err
	}

	return accessToken, refreshToken, nil
}

// VerifyToken checks jwt validity, expire date, blacklisted
func VerifyToken(token string, tokenType string) (*db.Token, error) {
	claims := &db.UserClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(Config.JWTSecretKey), nil
	})

	if err != nil || claims.Type != tokenType {
		return nil, errors.New("not valid token")
	}

	if time.Now().Sub(claims.ExpiresAt.Time) > 10*time.Second {
		return nil, errors.New("token is expired")
	}

	tokenModel := &db.Token{}
	userId := claims.Subject

	if err := DbConnection.Where("id = ? AND type >= ? AND blacklisted = ?", userId, tokenType, false).First(&tokenModel).Error; err != nil {
		return &db.Token{}, errors.New("cannot find token")
	}
	return tokenModel, nil
}
