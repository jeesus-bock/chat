package api

import (
	"chat/models"

	"github.com/gofiber/fiber/v2"
)

func getRoomsHandler(cfg *models.Config) func(c *fiber.Ctx) error {
	fn := func(c *fiber.Ctx) error {
		c.JSON(GetRooms(cfg, cc.sendConns))
		return nil
	}
	return fn
}
