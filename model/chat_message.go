package model

import "time"

type ChatMessage struct {
	Id       string
	Sender   *User
	ChatRoom *ChatRoom
	Message  string
	SendAt   time.Time
}

type ChatMessages []*ChatMessage

// for sorting need to implement go sort interface
func (slice ChatMessages) Len() int {
	return len(slice)
}

// sort ascending
func (slice ChatMessages) Less(i, j int) bool {
	return slice[i].SendAt.Unix() < slice[j].SendAt.Unix()
}

func (slice ChatMessages) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
