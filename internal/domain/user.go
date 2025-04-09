package domain

import "github.com/google/uuid"

type User struct {
	Id   uuid.UUID
	Name string
}

type DrawConnection struct {
	User     User
	OutputCh chan []Pixel
}

type ChatConnection struct {
	User     User
	OutputCh chan ChatMessage
}
