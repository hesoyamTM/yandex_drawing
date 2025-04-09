package domain

import "github.com/google/uuid"

type DrawEvent struct {
	UserId uuid.UUID
	Pixels []Pixel
}

type Pixel struct {
	Size  int    `json:"size"`
	Color string `json:"color"`
	X     int    `json:"x"`
	Y     int    `json:"y"`
}
