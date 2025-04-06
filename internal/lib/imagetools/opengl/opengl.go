package opengl

import (
	"image/color"
	"math"
	"runtime"

	"github.com/go-gl/gl/v4.5-compatibility/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/google/uuid"
	"github.com/hesoyamTM/yandex_drawing/internal/domain"
)

var (
	canvasStateCh chan CanvasState
)

type CanvasState struct {
	Id        uuid.UUID
	Width     int
	Height    int
	NewPixels []domain.Pixel
	Image     []byte
}

type UpdatedCanvasState struct {
	Id    uuid.UUID
	Image []byte
}

func StartOpenGL(maxW, maxH int) chan UpdatedCanvasState {
	canvasStateCh = make(chan CanvasState)
	resultStateCh := make(chan UpdatedCanvasState)

	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		defer close(canvasStateCh)
		defer close(resultStateCh)

		if err := glfw.Init(); err != nil {
			panic(err)
		}
		defer glfw.Terminate()

		glfw.WindowHint(glfw.Visible, glfw.False)
		window, err := glfw.CreateWindow(maxW, maxH, "", nil, nil)
		if err != nil {
			panic(err)
		}
		defer window.Destroy()

		window.MakeContextCurrent()

		if err := gl.Init(); err != nil {
			panic(err)
		}

		gl.MatrixMode(gl.PROJECTION)
		gl.LoadIdentity()
		gl.Ortho(0, float64(maxW), float64(maxH), 0, -1, 1)
		gl.MatrixMode(gl.MODELVIEW)
		gl.ClearColor(1, 1, 1, 1)

		for canvasState := range canvasStateCh {
			gl.Clear(gl.COLOR_BUFFER_BIT)

			if len(canvasState.Image) > 0 {
				gl.DrawPixels(
					int32(canvasState.Width),
					int32(canvasState.Height),
					gl.RGBA,
					gl.UNSIGNED_BYTE,
					gl.Ptr(canvasState.Image),
				)
			}

			for _, pix := range canvasState.NewPixels {
				drawCircle(pix, canvasState.Height)
			}

			pixels, _ := readPixels(canvasState.Width, canvasState.Height)
			resultStateCh <- UpdatedCanvasState{
				Id:    canvasState.Id,
				Image: pixels,
			}
		}
	}()

	return resultStateCh
}

func GetCanvasCh() chan<- CanvasState {
	return canvasStateCh
}

func drawCircle(pixel domain.Pixel, h int) {
	x := float32(pixel.X)
	y := float32(h - pixel.Y)
	radius := float32(pixel.Size/2) + 1

	color, _ := parseHexColorFast(pixel.Color)
	r := float32(color.R) / 255
	g := float32(color.G) / 255
	b := float32(color.B) / 255

	gl.Begin(gl.TRIANGLE_FAN)
	gl.Color3f(r, g, b)
	gl.Vertex2f(x, y) // Центр круга
	for i := 0; i <= 100; i++ {
		angle := 2 * math.Pi * float64(i) / 100
		x := x + float32(math.Cos(angle))*radius
		y := y + float32(math.Sin(angle))*radius
		gl.Vertex2f(x, y)
	}
	gl.End()
}

func readPixels(w, h int) ([]byte, error) {
	pixels := make([]byte, w*h*4)
	gl.ReadPixels(0, 0, int32(w), int32(h), gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels))
	return pixels, nil
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
