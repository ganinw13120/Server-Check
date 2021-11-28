package handler

import (
	"github.com/gofiber/fiber/v2"
)

func ParseRequest(c *fiber.Ctx, request interface{}) error {
	err := c.BodyParser(&request)
	if err != nil {
		return err
	}
	return nil
}
