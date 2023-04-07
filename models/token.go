package models

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

type UserClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
	Type  string `json:"type"`
}

type Token struct {
	ID          int64     `json:"id" gorm:"column:id;primary_key"`
	Token       string    `json:"token" bson:"token"`
	Type        string    `json:"type" bson:"type"`
	ExpiresAt   time.Time `json:"expires_at" bson:"expires_at"`
	Blacklisted bool      `json:"blacklisted" bson:"blacklisted"`
}

func (model Token) GetResponseJson() gin.H {
	return gin.H{"token": model.Token, "expires": model.ExpiresAt.Format("2006-01-02 15:04:05")}
}

func (Token) TableName() string {
	return "tokens"
}
