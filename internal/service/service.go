package service

import (
	"context"

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
}

func New(ctx context.Context, repo Repository) *DrawService {
	return &DrawService{
		repo: repo,
	}
}

func (d *DrawService) JoinToRoom(ctx context.Context, userId, canvasId int, inputCh <-chan []domain.Pixel) (<-chan []domain.Pixel, error) {
	if !d.repo.HasRoom(ctx, canvasId) {
		err := d.repo.CreateRoom(ctx, canvasId)
		if err != nil {
			// TODO: error
			return nil, err
		}
		d.workers[canvasId] = NewWorker(ctx, d.repo)
	}

	outputCh := make(chan []domain.Pixel, 100)

	if err := d.repo.JoinToRoom(ctx, canvasId, userId); err != nil {
		// TODO: error
		return nil, err
	}

	return outputCh, nil
}

func (d *DrawService) RemoveFromRoom(ctx context.Context, canvasId, userId int) error {
	if err := d.repo.RemoveFromRoom(ctx, canvasId, userId); err != nil {
		// TODO: error
		return err
	}

	return nil
}
