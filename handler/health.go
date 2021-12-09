package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/line/line-bot-sdk-go/v7/linebot"

	"server-health/model"
	"server-health/service"
)

type healthHandler struct {
	healthService service.IHealthService
	linebot       *linebot.Client
}

type IHealthHandler interface {
	WebHookHandler(*fiber.Ctx) error
}

func NewHealthHandler(service service.IHealthService, bot *linebot.Client) healthHandler {
	return healthHandler{
		healthService: service,
		linebot:       bot,
	}
}

func (h healthHandler) WebHookHandler(c *fiber.Ctx) error {
	request := model.LineWebhook{}
	err := ParseRequest(c, &request)
	if err != nil {
		fmt.Println(err)
		return c.Status(http.StatusOK).JSON(err)
	}
	var raw map[interface{}]interface{}
	err = ParseRequest(c, &raw)
	fmt.Println(err)
	fmt.Println(raw)
	raw2, err := json.Marshal(raw)
	fmt.Println(raw2)
	err = h.healthService.WebhookEnter(request)
	if err != nil {
		fmt.Println(err)
		return c.Status(http.StatusOK).JSON(err)
	}
	type test struct {
		gan string
	}
	return c.Status(http.StatusOK).JSON(map[string]string{
		"status": "success",
	})
}
