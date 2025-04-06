package inmemory

import (
	"context"
	"sync"

	"github.com/hesoyamTM/yandex_drawing/internal/domain"
)

type DrawRepository struct {
	rooms map[int]domain.Room
	l     sync.Mutex
}

func New() *DrawRepository {
	return &DrawRepository{
		rooms: make(map[int]domain.Room, 0),
	}
}

func (d *DrawRepository) CreateRoom(ctx context.Context, canvasId int) error {
	d.l.Lock()
	d.rooms[canvasId] = domain.Room{
		CanvasId:    canvasId,
		Canvas:      make([]domain.Pixel, 0),
		ActiveUsers: make([]int, 0),
	}
	d.l.Unlock()
	return nil
}

func (d *DrawRepository) HasRoom(ctx context.Context, canvasId int) bool {
	d.l.Lock()
	_, ok := d.rooms[canvasId]
	d.l.Unlock()

	return ok
}

func (d *DrawRepository) GetRoom(ctx context.Context, canvasId int) (domain.Room, error) {
	d.l.Lock()
	room, ok := d.rooms[canvasId]
	d.l.Unlock()

	if !ok {
		// TODO: error
		return domain.Room{}, nil
	}

	return room, nil
}

func (d *DrawRepository) GetRoomList(ctx context.Context) ([]domain.Room, error) {
	rooms := make([]domain.Room, 0, len(d.rooms))

	d.l.Lock()
	for _, room := range d.rooms {
		rooms = append(rooms, room)
	}
	d.l.Unlock()

	return rooms, nil
}

func (d *DrawRepository) JoinToRoom(ctx context.Context, canvasId int, userId int) (domain.Room, error) {
	d.l.Lock()
	room, ok := d.rooms[canvasId]

	if !ok {
		// TODO: error
		return domain.Room{}, nil
	}

	room.ActiveUsers = append(room.ActiveUsers, userId)
	d.rooms[canvasId] = room
	d.l.Unlock()

	return room, nil
}

func (d *DrawRepository) RemoveFromRoom(ctx context.Context, canvasId, userId int) error {
	d.l.Lock()
	room, ok := d.rooms[canvasId]

	if !ok {
		// TODO: error
		return nil
	}

	for i, user := range room.ActiveUsers {
		if user == userId {
			room.ActiveUsers[i] = room.ActiveUsers[len(room.ActiveUsers)-1]
			room.ActiveUsers = room.ActiveUsers[:len(room.ActiveUsers)-1]
			break
		}
	}

	d.rooms[canvasId] = room
	d.l.Unlock()

	return nil
}

func (d *DrawRepository) DeleteRoom(ctx context.Context, canvasId int) error {
	d.l.Lock()
	_, ok := d.rooms[canvasId]
	if !ok {
		// TODO: error
		return nil
	}

	delete(d.rooms, canvasId)
	d.l.Unlock()

	return nil
}
