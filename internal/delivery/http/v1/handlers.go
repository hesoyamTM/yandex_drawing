package v1

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/hesoyamTM/yandex_drawing/internal/domain"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

type DrawService interface {
	AddToCanvas(ctx context.Context, canvasId, userId uuid.UUID, username string, inputCh <-chan domain.DrawEvent) (<-chan []domain.Pixel, error)
	RemoveFromCanvas(ctx context.Context, canvasId, userId uuid.UUID) error
	AddToChat(ctx context.Context, canvasId, userId uuid.UUID, userName string, inputCh <-chan domain.ChatMessage) (<-chan domain.ChatMessage, error)
	RemoveFromChat(ctx context.Context, canvasId, userId uuid.UUID) error
	GetCanvas(ctx context.Context, canvasId, userId uuid.UUID) ([]byte, error)
}

func GetCanvas(drawService DrawService) echo.HandlerFunc {
	return func(c echo.Context) error {
		canvasId, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Logger().Error(err)
			return err
		}

		userId, err := uuid.Parse(c.FormValue("uid"))
		if err != nil {
			c.Logger().Error(err)
			return err
		}

		canvas, err := drawService.GetCanvas(c.Request().Context(), canvasId, userId)
		if err != nil {
			c.Logger().Error(err)
			return err
		}

		c.Blob(http.StatusOK, "image/png", canvas)
		return nil
	}
}

func Drawing(drawService DrawService) echo.HandlerFunc {
	return func(c echo.Context) error {
		canvasId, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Logger().Error(err)
			return err
		}

		userName := c.FormValue("name")
		userId, err := uuid.Parse(c.FormValue("uid"))
		if err != nil {
			c.Logger().Error(err)
			return err
		}

		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			c.Logger().Error(err)
			return err
		}
		defer ws.Close()

		inputCh := make(chan domain.DrawEvent)
		defer close(inputCh)

		outputCh, err := drawService.AddToCanvas(c.Request().Context(), userId, canvasId, userName, inputCh)
		if err != nil {
			c.Logger().Error(err)
			return err
		}
		defer drawService.RemoveFromCanvas(c.Request().Context(), canvasId, userId)

		go func() {
			for {
				pixels, ok := <-outputCh
				if !ok {
					ws.Close()
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

func Chat(drawService DrawService) echo.HandlerFunc {
	return func(c echo.Context) error {
		canvasId, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Logger().Error(err)
			return err
		}

		userName := c.FormValue("name")
		userId, err := uuid.Parse(c.FormValue("uid"))
		if err != nil {
			c.Logger().Error(err)
			return err
		}

		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			c.Logger().Error(err)
			return err
		}
		defer ws.Close()

		inputCh := make(chan domain.ChatMessage)
		defer close(inputCh)

		outputCh, err := drawService.AddToChat(c.Request().Context(), userId, canvasId, userName, inputCh)
		if err != nil {
			c.Logger().Error(err)
			return err
		}
		defer drawService.RemoveFromChat(c.Request().Context(), canvasId, userId)

		go func() {
			for {
				message, ok := <-outputCh
				if !ok {
					ws.Close()
					return
				}
				data, err := json.Marshal(message)
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

			var message domain.ChatMessage
			if err = json.Unmarshal(data, &message); err != nil {
				c.Logger().Error(err)
				return err
			}

			inputCh <- message
		}
	}
}
