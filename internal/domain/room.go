package domain

type Room struct {
	CanvasId    int
	Canvas      []Pixel
	ActiveUsers []int
}
