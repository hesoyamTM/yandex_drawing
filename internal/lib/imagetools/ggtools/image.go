package ggtools

import (
	"github.com/fogleman/gg"
	"github.com/hesoyamTM/yandex_drawing/internal/domain"
)

type Canvas struct {
	img *gg.Context
}

func New(w, h int) *Canvas {
	img := gg.NewContext(w, h)
	img.SetRGB(1, 1, 1)
	img.Clear()

	return &Canvas{
		img: img,
	}
}

func (c *Canvas) DrawCircles(pixels []domain.Pixel) {
	for _, pixel := range pixels {
		c.DrawCircle(pixel.X, pixel.Y, pixel.Size/2, pixel.Color)
	}
}

func (c *Canvas) DrawCircle(cx, cy, radius int, color string) {
	c.img.SetHexColor(color)
	c.img.DrawPoint(float64(cx), float64(cy), float64(radius))
}

func (c *Canvas) GetInBytes() []byte {
	img := c.img.Image()
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	pixels := make([]byte, w*h*4)
	offset := 0

	for y := range h {
		for x := range w {
			r, g, b, a := img.At(x, y).RGBA()
			pixels[offset] = byte(r >> 8)
			pixels[offset+1] = byte(g >> 8)
			pixels[offset+2] = byte(b >> 8)
			pixels[offset+3] = byte(a >> 8)
			offset += 4
		}
	}

	return pixels
}
