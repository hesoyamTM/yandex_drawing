package service

import (
	"context"
	"sync"

	"github.com/hesoyamTM/yandex_drawing/internal/domain"
)

type RoomWorker struct {
	InputCh chan domain.DrawEvent

	canvasId          int
	activeConnections map[int]*domain.Connection
	repo              Repository
	mut               sync.RWMutex
}

func NewWorker(ctx context.Context, canvasId int, repo Repository) (*RoomWorker, error) {
	// err := repo.CreateRoom(ctx, canvasId)
	// if err != nil {
	// 	// TODO: error
	// 	return nil, err
	// }

	worker := &RoomWorker{
		InputCh:           make(chan domain.DrawEvent, 100),
		activeConnections: make(map[int]*domain.Connection, 0),
		repo:              repo,
	}

	return worker, nil
}

func (r *RoomWorker) AddConnection(ctx context.Context, conn *domain.Connection) error {
	r.mut.Lock()
	r.activeConnections[conn.UserId] = conn
	r.mut.Unlock()

	return nil
}

func (r *RoomWorker) RemoveConnection(ctx context.Context, userId int) error {
	r.mut.Lock()
	close(r.activeConnections[userId].OutputCh)
	delete(r.activeConnections, userId)
	r.mut.Unlock()

	return nil
}

func (r *RoomWorker) Stop(ctx context.Context) {
	close(r.InputCh)
}

func (r *RoomWorker) Run(ctx context.Context) {
	for {
		event, ok := <-r.InputCh
		if !ok {
			return
		}

		r.mut.RLock()
		for _, conn := range r.activeConnections {
			if conn.UserId != event.UserId {
				conn.OutputCh <- event.Pixels
			}
		}
		r.mut.RUnlock()
	}
}
