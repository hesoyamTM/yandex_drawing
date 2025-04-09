package canvas

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/hesoyamTM/yandex_drawing/internal/domain"
)

type ChatBroadcaster struct {
	InputCh chan domain.ChatMessage

	activeConnections map[uuid.UUID]*domain.ChatConnection

	l sync.Mutex
}

func NewChatBroadcaster(ctx context.Context) *ChatBroadcaster {
	return &ChatBroadcaster{
		InputCh:           make(chan domain.ChatMessage, 10000),
		activeConnections: make(map[uuid.UUID]*domain.ChatConnection, 0),
	}
}

func (c *ChatBroadcaster) AddConnection(ctx context.Context, userId uuid.UUID, conn *domain.ChatConnection) {
	c.l.Lock()
	defer c.l.Unlock()

	c.activeConnections[userId] = conn
}

func (c *ChatBroadcaster) RemoveConnection(ctx context.Context, userId uuid.UUID) {
	c.l.Lock()
	defer c.l.Unlock()

	close(c.activeConnections[userId].OutputCh)
	delete(c.activeConnections, userId)
}

func (c *ChatBroadcaster) Run(ctx context.Context) {
	for message := range c.InputCh {
		c.l.Lock()

		for _, conn := range c.activeConnections {
			conn.OutputCh <- message
		}
		c.l.Unlock()
	}
}

func (c *ChatBroadcaster) Stop() {
	close(c.InputCh)
}
