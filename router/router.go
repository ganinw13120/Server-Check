package router

import (
	"github.com/gofiber/fiber/v2"

	"server-health/handler"
)

type route struct {
	app *fiber.App

	healthHandler handler.IHealthHandler
}

var router *route = nil

func New(app *fiber.App, healthHandler handler.IHealthHandler) *route {
	if router == nil {
		router = &route{
			app:           app,
			healthHandler: healthHandler,
		}
		router.setUp()
	}
	return router
}

func (r route) setUp() {
	group := r.app.Group("webhook")
	{
		group.Post("/", r.healthHandler.WebHookHandler)
	}
}
