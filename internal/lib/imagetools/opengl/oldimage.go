package opengl

// import (
// 	"errors"
// 	"image/color"
// 	"math"
// 	"runtime"

// 	"github.com/go-gl/gl/v4.5-compatibility/gl"
// 	"github.com/go-gl/glfw/v3.2/glfw"
// 	"github.com/hesoyamTM/yandex_drawing/internal/domain"
// )

// var errInvalidFormat = errors.New("invalid format")

// type Canvas struct {
// 	pixels []domain.Pixel
// 	width  int
// 	height int
// }

// func New(w, h int) *Canvas {
// 	return &Canvas{
// 		pixels: make([]domain.Pixel, 0),
// 		width:  w,
// 		height: h,
// 	}
// }

// func (c *Canvas) DrawCircles(pixels []domain.Pixel) {
// 	c.pixels = append(c.pixels, pixels...)
// }

// func (c *Canvas) DrawCircle(cx, cy, radius int, color string) {

// }

// func (c *Canvas) GetInBytes() []byte {
// 	runtime.LockOSThread()
// 	defer runtime.UnlockOSThread()

// 	if err := glfw.Init(); err != nil {
// 		panic(err)
// 	}
// 	defer glfw.Terminate()

// 	glfw.WindowHint(glfw.Visible, glfw.False)
// 	window, err := glfw.CreateWindow(c.width, c.height, "", nil, nil)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer window.Destroy()

// 	window.MakeContextCurrent()

// 	if err := gl.Init(); err != nil {
// 		panic(err)
// 	}

// 	gl.MatrixMode(gl.PROJECTION)
// 	gl.LoadIdentity()
// 	gl.Ortho(0, float64(c.width), float64(c.height), 0, -1, 1)
// 	gl.MatrixMode(gl.MODELVIEW)
// 	gl.ClearColor(1, 1, 1, 1)

// 	gl.Clear(gl.COLOR_BUFFER_BIT)

// 	for _, pix := range c.pixels {
// 		c.drawCircle(pix)
// 	}

// 	pixels, _ := readPixels(c.width, c.height)

// 	return pixels
// }

// func (c *Canvas) drawCircle(pixel domain.Pixel) {
// 	x := float32(pixel.X)
// 	y := float32(c.height - pixel.Y)
// 	radius := float32(pixel.Size / 2)

// 	color, _ := parseHexColorFast(pixel.Color)
// 	r := float32(color.R) / 255
// 	g := float32(color.G) / 255
// 	b := float32(color.B) / 255

// 	gl.Begin(gl.TRIANGLE_FAN)
// 	gl.Color3f(r, g, b)
// 	gl.Vertex2f(x, y) // Центр круга
// 	for i := 0; i <= 100; i++ {
// 		angle := 2 * math.Pi * float64(i) / 100
// 		x := x + float32(math.Cos(angle))*radius
// 		y := y + float32(math.Sin(angle))*radius
// 		gl.Vertex2f(x, y)
// 	}
// 	gl.End()
// }

// func parseHexColorFast(s string) (c color.RGBA, err error) {
// 	c.A = 0xff

// 	if s[0] != '#' {
// 		return c, errInvalidFormat
// 	}

// 	hexToByte := func(b byte) byte {
// 		switch {
// 		case b >= '0' && b <= '9':
// 			return b - '0'
// 		case b >= 'a' && b <= 'f':
// 			return b - 'a' + 10
// 		case b >= 'A' && b <= 'F':
// 			return b - 'A' + 10
// 		}
// 		err = errInvalidFormat
// 		return 0
// 	}

// 	switch len(s) {
// 	case 7:
// 		c.R = hexToByte(s[1])<<4 + hexToByte(s[2])
// 		c.G = hexToByte(s[3])<<4 + hexToByte(s[4])
// 		c.B = hexToByte(s[5])<<4 + hexToByte(s[6])
// 	case 4:
// 		c.R = hexToByte(s[1]) * 17
// 		c.G = hexToByte(s[2]) * 17
// 		c.B = hexToByte(s[3]) * 17
// 	default:
// 		err = errInvalidFormat
// 	}
// 	return
// }
