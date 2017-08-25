package model

import "time"

type ChatRoom struct {
	Id           string
	Name         string
	Participants []*User
	CreatedAt    time.Time
}
