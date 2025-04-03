package domain

type DrawEvent struct {
	UserId int
	Pixels []Pixel
}

type Pixel struct {
	Size  float32 `json:"size"`
	Color string  `json:"color"`
	X     float32 `json:"x"`
	Y     float32 `json:"y"`
}
