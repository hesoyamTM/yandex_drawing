package service

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/hesoyamTM/yandex_drawing/internal/domain"
)

const (
	channelBuffer = 10
)

type Repository interface {
	CreateRoom(ctx context.Context, canvasId int) error
	HasRoom(ctx context.Context, canvasId int) bool
	GetRoom(ctx context.Context, canvasId int) (domain.Room, error)
	GetRoomList(ctx context.Context) ([]domain.Room, error)
	JoinToRoom(ctx context.Context, canvasId int, userId int) (domain.Room, error)
	RemoveFromRoom(ctx context.Context, canvasId, userId int) error
	DeleteRoom(ctx context.Context, canvasId int) error
}

type DrawService struct {
	repo    Repository
	workers map[uuid.UUID]*RoomWorker
	l       sync.Mutex
}

func New(ctx context.Context, repo Repository) *DrawService {
	return &DrawService{
		repo:    repo,
		workers: make(map[uuid.UUID]*RoomWorker),
	}
}

func (d *DrawService) GetCanvas(ctx context.Context, canvasId, userId uuid.UUID) ([]byte, error) {

	d.l.Lock()
	defer d.l.Unlock()
	worker, ok := d.workers[canvasId]

	if !ok {
		worker, err := NewWorker(ctx, canvasId, d.repo)
		if err != nil {
			// TODO: error
			return nil, err
		}

		go worker.CanvasBroadcaster.Run(ctx)

		d.workers[canvasId] = worker
		return worker.CanvasBroadcaster.GetCanvas(ctx, userId)

		// TODO: если нет подключений, то убираем воркер
	}

	return worker.CanvasBroadcaster.GetCanvas(ctx, userId)
}

func (d *DrawService) AddToCanvas(ctx context.Context, canvasId, userId uuid.UUID, username string, inputCh <-chan domain.DrawEvent) (<-chan []domain.Pixel, error) {
	_, ok := d.workers[canvasId]

	if !ok {
		worker, err := NewWorker(ctx, canvasId, d.repo)
		if err != nil {
			// TODO: error
			return nil, err
		}

		go worker.CanvasBroadcaster.Run(ctx)

		d.l.Lock()
		d.workers[canvasId] = worker
		d.l.Unlock()
	}

	outputCh := make(chan []domain.Pixel, channelBuffer)
	conn := &domain.DrawConnection{
		User: domain.User{
			Id:   userId,
			Name: "",
		},
		OutputCh: outputCh,
	}

	d.l.Lock()
	d.workers[canvasId].CanvasBroadcaster.AddConnection(ctx, conn)
	broadcastCh := d.workers[canvasId].CanvasBroadcaster.InputCh
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

func (d *DrawService) RemoveFromCanvas(ctx context.Context, canvasId, userId uuid.UUID) error {
	d.l.Lock()
	worker := d.workers[canvasId]
	d.l.Unlock()

	worker.CanvasBroadcaster.RemoveConnection(ctx, userId)

	if worker.CanvasBroadcaster.GetActiveConnectionCount() == 0 {
		worker.CanvasBroadcaster.Stop()

		d.l.Lock()
		delete(d.workers, canvasId)
		d.l.Unlock()
	}

	return nil
}

func (d *DrawService) AddToChat(ctx context.Context, canvasId, userId uuid.UUID, username string, inputCh <-chan domain.ChatMessage) (<-chan domain.ChatMessage, error) {
	return nil, nil
}

func (d *DrawService) RemoveFromChat(ctx context.Context, canvasId, userId uuid.UUID) error {
	return nil
}

func (d *DrawService) addWorker(ctx context.Context, canvasId uuid.UUID) error {
	worker, err := NewWorker(ctx, canvasId, d.repo)
	if err != nil {
		// TODO: error
		return err
	}

	go worker.CanvasBroadcaster.Run(ctx)

	d.l.Lock()
	d.workers[canvasId] = worker
	d.l.Unlock()

	return nil
}
