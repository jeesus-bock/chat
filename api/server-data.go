package api

import (
	"chat/models"

	"github.com/gofiber/fiber/v2"
)

func serverDataHandler(c *fiber.Ctx) error {
	s := new(models.ServerData)
	s.Rooms = make([]models.Room, 0)
	s.Users = make([]models.User, 0)

	c.JSON(cc.sendConns)
	return nil
}
