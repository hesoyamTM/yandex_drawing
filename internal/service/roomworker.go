package service

import (
	"context"
	"sync"

	"github.com/hesoyamTM/yandex_drawing/internal/domain"
)

type RoomWorker struct {
	InputCh           chan domain.DrawEvent
	activeConnections []domain.Connection
	repo              Repository
	mut               sync.Mutex
}

func NewWorker(ctx context.Context, repo Repository) *RoomWorker {
	worker := &RoomWorker{
		InputCh:           make(chan domain.DrawEvent, 100),
		activeConnections: make([]domain.Connection, 0),
		repo:              repo,
	}

	go worker.run(ctx)

	return worker
}

func (r *RoomWorker) AddConnection() error {

}

func (r *RoomWorker) RemoveConnection() error {

}

func (r *RoomWorker) run(ctx context.Context) {
	for {
		event, ok := <-r.InputCh
		if !ok {
			return
		}

		r.mut.Lock()
		for _, conn := range r.activeConnections {
			if conn.UserId != event.UserId {
				conn.OutputCh <- event.Pixels
			}
		}
		r.mut.Lock()
	}
}
