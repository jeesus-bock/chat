package ws

import "github.com/gofiber/fiber/v2"

func monitor(c *fiber.Ctx) error {
	log.Info("cc: ", cc.sendConns)
	c.JSON(cc.sendConns)
	return nil
}
