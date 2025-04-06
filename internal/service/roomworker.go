package service

import (
	"context"
	"sync"

	"github.com/hesoyamTM/yandex_drawing/internal/domain"
	"github.com/hesoyamTM/yandex_drawing/internal/lib/imagetools/opengl"
)

const (
	h = 720
	w = 1280
)

type Canvas interface {
	DrawCircles(pixels []domain.Pixel)
	DrawCircle(cx, cy, radius int, color string)
	GetInBytes() []byte
}

type RoomWorker struct {
	InputCh chan domain.DrawEvent

	canvasId          int
	activeConnections map[int]*domain.Connection
	repo              Repository

	canvas Canvas

	mut sync.Mutex
}

func NewWorker(ctx context.Context, canvasId int, repo Repository) (*RoomWorker, error) {
	worker := &RoomWorker{
		canvasId:          canvasId,
		InputCh:           make(chan domain.DrawEvent, 10000),
		activeConnections: make(map[int]*domain.Connection, 0),
		repo:              repo,
		canvas:            opengl.New(w, h),
	}

	repo.CreateRoom(ctx, canvasId)

	return worker, nil
}

func (r *RoomWorker) LockBroadcast() {
	r.mut.Lock()
}

func (r *RoomWorker) UnlockBroadcast() {
	r.mut.Unlock()
}

func (r *RoomWorker) GetCanvas(ctx context.Context) Canvas {
	return r.canvas
}

func (r *RoomWorker) AddConnection(ctx context.Context, conn *domain.Connection) error {
	r.mut.Lock()
	defer r.mut.Unlock()

	r.activeConnections[conn.UserId] = conn

	_, err := r.repo.JoinToRoom(ctx, r.canvasId, conn.UserId)
	if err != nil {
		// TODO: error
		return err
	}

	return nil
}

func (r *RoomWorker) RemoveConnection(ctx context.Context, userId int) error {
	r.mut.Lock()
	defer r.mut.Unlock()

	close(r.activeConnections[userId].OutputCh)
	delete(r.activeConnections, userId)

	if err := r.repo.RemoveFromRoom(ctx, r.canvasId, userId); err != nil {
		// TODO: error
		return err
	}

	return nil
}

func (r *RoomWorker) Stop(ctx context.Context) {
	close(r.InputCh)
}

func (r *RoomWorker) Run(ctx context.Context) {
	for event := range r.InputCh {
		r.mut.Lock()
		r.canvas.DrawCircles(event.Pixels)

		for _, conn := range r.activeConnections {
			if conn.UserId != event.UserId {
				conn.OutputCh <- event.Pixels
			}
		}
		r.mut.Unlock()
	}
}
