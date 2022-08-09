package api

import (
	"chat/models"

	"github.com/gofiber/fiber/v2"
)

func getServerDataHandler(cfg *models.Config) func(c *fiber.Ctx) error {
	fn := func(c *fiber.Ctx) error {
		s := new(models.ServerData)
		s.Name = cfg.Name
		s.Type = cfg.Type
		s.URL = cfg.Host
		s.VoiceURL = cfg.Host
		c.JSON(*s)
		return nil
	}
	return fn
}

func GetRoomsMapHandler(c *fiber.Ctx) error {
	return c.JSON(GetRoomsMap())
}
