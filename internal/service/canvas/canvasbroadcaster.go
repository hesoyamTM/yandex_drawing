package canvas

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hesoyamTM/yandex_drawing/internal/domain"
	"github.com/hesoyamTM/yandex_drawing/internal/lib/imagetools/opengl"
)

const (
	timeToConnect = 5 * time.Second
)

type Canvas interface {
	DrawCircles(pixels []domain.Pixel)
	GetInBytes() []byte
}

type CanvasBroadcaster struct {
	InputCh chan domain.DrawEvent

	activeConnections  map[uuid.UUID]*domain.DrawConnection
	waitingConnections map[uuid.UUID]*domain.WaitingConnection
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

func (c *CanvasBroadcaster) GetActiveConnectionCount() int {
	c.l.Lock()
	defer c.l.Unlock()
	return len(c.activeConnections)
}

func (c *CanvasBroadcaster) GetCanvas(ctx context.Context, userId uuid.UUID) ([]byte, error) {
	c.l.Lock()
	defer c.l.Unlock()

	c.waitingConnections[userId] = &domain.WaitingConnection{
		UserId:      userId,
		Connected:   false,
		PixelBuffer: make([]domain.Pixel, 0),
	}

	go func() {
		// Если через определённое время пользователь не подключился, то мы удаляем весь его буфер
		<-time.After(timeToConnect)
		c.l.Lock()
		defer c.l.Unlock()

		if !c.waitingConnections[userId].Connected {
			delete(c.waitingConnections, userId)
		}
	}()

	return c.canvas.GetInBytes(), nil
}

func (c *CanvasBroadcaster) AddConnection(ctx context.Context, conn *domain.DrawConnection) {
	c.l.Lock()
	defer c.l.Unlock()

	c.activeConnections[conn.User.Id] = conn
	wConn, ok := c.waitingConnections[conn.User.Id]

	if ok {
		conn.OutputCh <- wConn.PixelBuffer
		delete(c.waitingConnections, conn.User.Id)
	}
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

		for _, wCoon := range c.waitingConnections {
			wCoon.PixelBuffer = append(wCoon.PixelBuffer, event.Pixels...)
		}

		for _, conn := range c.activeConnections {
			conn.OutputCh <- event.Pixels
		}
		c.l.Unlock()
	}
}

func (c *CanvasBroadcaster) Stop() {
	close(c.InputCh)
}
