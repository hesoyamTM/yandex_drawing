package domain

type Connection struct {
	UserId   int
	OutputCh chan []Pixel
}
