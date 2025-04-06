package opengl

import (
	"errors"

	"github.com/hesoyamTM/yandex_drawing/internal/domain"
)

var errInvalidFormat = errors.New("invalid format")

type Canvas struct {
	pixels       []domain.Pixel
	currentImage []byte
	width        int
	height       int
}

func New(w, h int) *Canvas {
	return &Canvas{
		pixels:       make([]domain.Pixel, 0),
		currentImage: make([]byte, 0),
		width:        w,
		height:       h,
	}
}

func (c *Canvas) DrawCircles(pixels []domain.Pixel) {
	c.pixels = append(c.pixels, pixels...)
	c.currentImage = <-UpdateCanvas(c.width, c.height, c.pixels, c.currentImage)
	c.pixels = c.pixels[:0]
}

func (c *Canvas) DrawCircle(cx, cy, radius int, color string) {

}

func (c *Canvas) GetInBytes() []byte {
	if len(c.pixels) == 0 {
		return c.currentImage
	}

	c.currentImage = <-UpdateCanvas(c.width, c.height, c.pixels, c.currentImage)
	c.pixels = c.pixels[:0]

	return c.currentImage
}
