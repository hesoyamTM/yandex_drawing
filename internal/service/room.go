package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/hesoyamTM/yandex_drawing/internal/service/canvas"
	"github.com/hesoyamTM/yandex_drawing/internal/service/chat"
)

const (
	h             = 720
	w             = 1280
	lastPixelSize = 100000
)

type RoomWorker struct {
	CanvasBroadcaster *canvas.CanvasBroadcaster
	ChatBroadcaster   *chat.ChatBroadcaster

	canvasId uuid.UUID
	repo     Repository
}

func NewWorker(ctx context.Context, canvasId uuid.UUID, repo Repository) (*RoomWorker, error) {
	// TODO: request to Canvas Service: getting width, heigh and other information

	worker := &RoomWorker{
		canvasId:          canvasId,
		repo:              repo,
		CanvasBroadcaster: canvas.NewCanvasBroadcaster(ctx, w, h),
		ChatBroadcaster:   chat.NewChatBroadcaster(ctx),
	}

	return worker, nil
}
