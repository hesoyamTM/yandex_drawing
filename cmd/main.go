package main

import (
	"context"

	v1 "github.com/hesoyamTM/yandex_drawing/internal/delivery/http/v1"
	"github.com/hesoyamTM/yandex_drawing/internal/lib/imagetools/opengl"
	"github.com/hesoyamTM/yandex_drawing/internal/repository/inmemory"
	"github.com/hesoyamTM/yandex_drawing/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	opengl.InitOpenGlScheduler(opengl.StartOpenGL(1280, 720))

	drawRepo := inmemory.New()
	drawService := service.New(context.Background(), drawRepo)

	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Static("/", "../public")
	e.GET("/ws/drawing", v1.HandleConnections(drawService))
	e.Logger.Fatal(e.Start(":1323"))
}
