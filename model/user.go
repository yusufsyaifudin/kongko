package model

import "time"

type User struct {
	Id           string
	Email        string
	Password     string
	RegisteredAt time.Time
}
