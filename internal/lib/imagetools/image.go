package imagetools

import (
	"errors"
	"image"
	"image/color"

	"github.com/hesoyamTM/yandex_drawing/internal/domain"
)

var errInvalidFormat = errors.New("invalid format")

type Canvas struct {
	img *image.RGBA
}

func New(w, h int) *Canvas {
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	for x := range w {
		for y := range h {
			img.Set(x, y, color.White)
		}
	}

	return &Canvas{
		img: img,
	}
}

func (c *Canvas) DrawCircles(pixels []domain.Pixel) {
	for _, pixel := range pixels {
		c.DrawCircle(pixel.X, pixel.Y, pixel.Size/2, pixel.Color)
	}
}

func (c *Canvas) DrawCircle(cx, cy, radius int, hexColor string) {
	color, _ := parseHexColorFast(hexColor)

	radiusSquared := radius * radius
	for y := cy - radius; y <= cy+radius; y++ {
		for x := cx - radius; x <= cx+radius; x++ {
			dx, dy := x-cx, y-cy
			if dx*dx+dy*dy <= radiusSquared && c.img.At(x, y) != color {
				c.img.Set(x, y, color)
			}
		}
	}
}

func (c *Canvas) GetInBytes() []byte {
	w, h := c.img.Rect.Dx(), c.img.Rect.Dy()
	pixels := make([]byte, w*h*4)
	offset := 0

	for y := range h {
		for x := range w {
			r, g, b, a := c.img.At(x, y).RGBA()
			pixels[offset] = byte(r >> 8)
			pixels[offset+1] = byte(g >> 8)
			pixels[offset+2] = byte(b >> 8)
			pixels[offset+3] = byte(a >> 8)
			offset += 4
		}
	}

	return pixels
}

func parseHexColorFast(s string) (c color.RGBA, err error) {
	c.A = 0xff

	if s[0] != '#' {
		return c, errInvalidFormat
	}

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		err = errInvalidFormat
		return 0
	}

	switch len(s) {
	case 7:
		c.R = hexToByte(s[1])<<4 + hexToByte(s[2])
		c.G = hexToByte(s[3])<<4 + hexToByte(s[4])
		c.B = hexToByte(s[5])<<4 + hexToByte(s[6])
	case 4:
		c.R = hexToByte(s[1]) * 17
		c.G = hexToByte(s[2]) * 17
		c.B = hexToByte(s[3]) * 17
	default:
		err = errInvalidFormat
	}
	return
}
