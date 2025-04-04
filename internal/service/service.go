package service

import (
	"context"
	"sync"

	"github.com/hesoyamTM/yandex_drawing/internal/domain"
)

type Repository interface {
	CreateRoom(ctx context.Context, canvasId int) error
	HasRoom(ctx context.Context, canvasId int) bool
	GetRoom(ctx context.Context, canvasId int) (domain.Room, error)
	GetRoomList(ctx context.Context) ([]domain.Room, error)
	JoinToRoom(ctx context.Context, canvasId int, userId int) error
	RemoveFromRoom(ctx context.Context, canvasId, userId int) error
	DeleteRoom(ctx context.Context, canvasId int) error
}

type DrawService struct {
	repo    Repository
	workers map[int]*RoomWorker
	l       sync.Mutex
}

func New(ctx context.Context, repo Repository) *DrawService {
	return &DrawService{
		repo:    repo,
		workers: make(map[int]*RoomWorker),
	}
}

func (d *DrawService) JoinToRoom(ctx context.Context, userId, canvasId int, inputCh <-chan domain.DrawEvent) (<-chan []domain.Pixel, error) {
	if len(d.workers) == 0 {
		worker, err := NewWorker(ctx, canvasId, d.repo)
		if err != nil {
			// TODO: error
			return nil, err
		}

		go worker.Run(ctx)

		d.l.Lock()
		d.workers[canvasId] = worker
		d.l.Unlock()
	}

	outputCh := make(chan []domain.Pixel, 100)
	conn := &domain.Connection{
		UserId:   userId,
		OutputCh: outputCh,
	}

	d.l.Lock()
	d.workers[canvasId].AddConnection(ctx, conn)
	broadcastCh := d.workers[canvasId].InputCh
	d.l.Unlock()

	go func() {
		for {
			event, ok := <-inputCh
			if !ok {
				return
			}

			broadcastCh <- event
		}
	}()

	return outputCh, nil
}

func (d *DrawService) RemoveFromRoom(ctx context.Context, canvasId, userId int) error {
	d.l.Lock()
	d.workers[canvasId].RemoveConnection(ctx, userId)
	d.l.Unlock()

	return nil
}
