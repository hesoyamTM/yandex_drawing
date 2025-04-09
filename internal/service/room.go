package service

import (
	"context"

	"github.com/google/uuid"
)

const (
	h             = 720
	w             = 1280
	lastPixelSize = 100000
)

type RoomWorker struct {
	CanvasBroadcaster *canvas.CanvasBroadcaster
	ChatBroadcaster   *canvas.ChatBroadcaster

	canvasId uuid.UUID
	repo     Repository
}

func NewWorker(ctx context.Context, canvasId uuid.UUID, repo Repository) (*RoomWorker, error) {
	// TODO: request to Canvas Service: getting width, heigh and other information

	worker := &RoomWorker{
		canvasId:          canvasId,
		repo:              repo,
		CanvasBroadcaster: NewCanvasBroadcaster(ctx, w, h),
		ChatBroadcaster:   NewChatBroadcaster(ctx),
	}

	return worker, nil
}
