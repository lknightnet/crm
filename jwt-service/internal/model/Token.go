package model

import "time"

type AccessToken struct {
	Issuer       string    //издатель
	Audience     string    //аудитория
	Subject      string    //объект
	ExpirationAt time.Time //время истечения
	IssuedAt     time.Time //время выпуска
	Token        string
}

type RefreshToken struct {
	Issuer       string    //издатель
	Audience     string    //аудитория
	Subject      string    //объект
	ExpirationAt time.Time //время истечения
	IssuedAt     time.Time //время выпуска
	Token        string
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
