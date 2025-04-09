package canvas

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/hesoyamTM/yandex_drawing/internal/domain"
	"github.com/hesoyamTM/yandex_drawing/internal/lib/imagetools/opengl"
)

type Canvas interface {
	DrawCircles(pixels []domain.Pixel)
	GetInBytes() []byte
}

type CanvasBroadcaster struct {
	InputCh chan domain.DrawEvent

	activeConnections  map[uuid.UUID]*domain.DrawConnection
	waitingConnections map[uuid.UUID]*domain.DrawConnection
	canvas             Canvas

	l sync.Mutex
}

func NewCanvasBroadcaster(ctx context.Context, w, h int) *CanvasBroadcaster {
	return &CanvasBroadcaster{
		InputCh:           make(chan domain.DrawEvent, 10000),
		activeConnections: make(map[uuid.UUID]*domain.DrawConnection, 0),
		canvas:            opengl.New(w, h),
	}
}

func (c *CanvasBroadcaster) GetCanvas(ctx context.Context, userId uuid.UUID, conn *domain.DrawConnection) {

}

func (c *CanvasBroadcaster) AddConnection(ctx context.Context, userId uuid.UUID, conn *domain.DrawConnection) {
	c.l.Lock()
	defer c.l.Unlock()

	c.activeConnections[userId] = conn
}

func (c *CanvasBroadcaster) RemoveConnection(ctx context.Context, userId uuid.UUID) {
	c.l.Lock()
	defer c.l.Unlock()

	close(c.activeConnections[userId].OutputCh)
	delete(c.activeConnections, userId)
}

func (c *CanvasBroadcaster) Run(ctx context.Context) {
	for event := range c.InputCh {
		c.l.Lock()
		c.canvas.DrawCircles(event.Pixels)

		for _, conn := range c.activeConnections {
			conn.OutputCh <- event.Pixels
		}
		c.l.Unlock()
	}
}

func (c *CanvasBroadcaster) Stop() {
	close(c.InputCh)
}
