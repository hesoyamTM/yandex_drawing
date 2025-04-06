package opengl

import (
	"github.com/google/uuid"
	"github.com/hesoyamTM/yandex_drawing/internal/domain"
)

var (
	waitingCanvases map[uuid.UUID]chan []byte
)

func InitOpenGlScheduler(inputCh <-chan UpdatedCanvasState) {
	waitingCanvases = make(map[uuid.UUID]chan []byte)

	go func() {
		for updCanvas := range inputCh {
			ch := waitingCanvases[updCanvas.Id]
			ch <- updCanvas.Image

			close(ch)
		}
	}()
}

func UpdateCanvas(w, h int, newPixels []domain.Pixel, image []byte) <-chan []byte {
	resCh := make(chan []byte)

	id := uuid.New()
	canvasState := CanvasState{
		Id:        id,
		NewPixels: newPixels,
		Width:     w,
		Height:    h,
		Image:     image,
	}

	waitingCanvases[id] = resCh
	GetCanvasCh() <- canvasState

	return resCh
}
