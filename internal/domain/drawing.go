package domain

type DrawEvent struct {
	UserId int
	Pixels []Pixel
}

type Pixel struct {
	Color Color `json:"color"`
	X     int   `json:"x"`
	Y     int   `json:"y"`
}

type Color struct {
	R uint8 `json:"R"`
	G uint8 `json:"G"`
	B uint8 `json:"B"`
}
