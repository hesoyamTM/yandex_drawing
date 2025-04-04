package v1

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/hesoyamTM/yandex_drawing/internal/domain"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	canvasId = 1
)

type DrawService interface {
	JoinToRoom(ctx context.Context, userId, canvasId int, inputCh <-chan domain.DrawEvent) (<-chan []domain.Pixel, error)
	RemoveFromRoom(ctx context.Context, canvasId, userId int) error
}

func HandleConnections(drawService DrawService) echo.HandlerFunc {
	return func(c echo.Context) error {
		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			c.Logger().Error(err)
			return err
		}
		defer ws.Close()

		inputCh := make(chan domain.DrawEvent, 100)
		defer close(inputCh)

		userId := rand.Int()

		outputCh, err := drawService.JoinToRoom(c.Request().Context(), userId, canvasId, inputCh)
		if err != nil {
			c.Logger().Error(err)
			return err
		}
		defer drawService.RemoveFromRoom(c.Request().Context(), canvasId, userId)

		go func() {
			for {
				pixels, ok := <-outputCh
				if !ok {
					return
				}
				data, err := json.Marshal(pixels)
				if err != nil {
					c.Logger().Error(err)
					return
				}
				if err = ws.WriteMessage(websocket.TextMessage, data); err != nil {
					c.Logger().Error(err)
					return
				}
			}
		}()

		for {
			_, data, err := ws.ReadMessage()
			if err != nil {
				c.Logger().Error(err)
				return err
			}
			var changedPixels []domain.Pixel
			if err = json.Unmarshal(data, &changedPixels); err != nil {
				c.Logger().Error(err)
				return err
			}
			inputCh <- domain.DrawEvent{
				UserId: userId,
				Pixels: changedPixels,
			}
		}
	}
}
