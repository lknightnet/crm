package model

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"time"
)

type RouteLink struct {
	OneTimeToken string
	CreatedAt    time.Time
	ExpirationAt time.Time
	UUID         string
}

func GenerateRouteLink(userUUID, email string) *RouteLink {
	createdAt := time.Now()
	return &RouteLink{
		OneTimeToken: generateToken(userUUID, email, createdAt),
		CreatedAt:    createdAt,
		ExpirationAt: time.Now().Add(time.Hour * 24),
		UUID:         userUUID,
	}
}

func generateToken(userUUID, email string, createdAt time.Time) string {
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Println(err)
	}
	tokenData := fmt.Sprintf("%s:%d:%x:%s", userUUID, createdAt.Unix(), randomBytes, email)

	token := base64.URLEncoding.EncodeToString([]byte(tokenData))

	return token
}
