package domain

type DrawEvent struct {
	UserId int
	Pixels []Pixel
}

type Pixel struct {
	Size  int    `json:"size"`
	Color string `json:"color"`
	X     int    `json:"x"`
	Y     int    `json:"y"`
}
