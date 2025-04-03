package main

import (
	"context"

	v1 "github.com/hesoyamTM/yandex_drawing/internal/delivery/http/v1"
	"github.com/hesoyamTM/yandex_drawing/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	drawService := service.New(context.Background(), nil)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Static("/", "../public")
	e.GET("/ws/drawing", v1.HandleConnections(drawService))
	e.Logger.Fatal(e.Start(":1323"))
}
