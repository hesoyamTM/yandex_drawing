package domain

const (
	xSize = 1280
	ySize = 720
)

type Room struct {
	CanvasId    int
	Canvas      [ySize][xSize]Color
	ActiveUsers []int
}
