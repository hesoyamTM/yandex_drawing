package domain

const (
	xSize = 1280
	ySize = 720
)

type Room struct {
	Canvas            [ySize][xSize]Color
	ActiveConnections map[int]*Connection
}

type Connection struct {
	InputCh  <-chan []Pixel
	OutputCh chan []Pixel
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
