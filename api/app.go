package api

import (
	"chat/logger"
	"chat/models"
	"fmt"
	"mime/multipart"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger
var app *fiber.App

func Init(cfg *models.Config) {
	log = logger.GetLogger()

	app = fiber.New()

	app.Use(cors.New(cors.Config{AllowOrigins: "*"}))
	app.Get("/server", getServerDataHandler(cfg))
	app.Get("/rooms", GetRoomsMapHandler)
	app.Static("/uploads", "./uploads")
	app.Post("/upload/:id/:room", func(c *fiber.Ctx) (err error) {
		log.Info("/upload POST handler")
		room := c.Params("room")
		id := c.Params("id")
		if room == "" || id == "" {
			c.SendStatus(fiber.StatusBadRequest)
			return
		}
		nick := cc.sendConns[id].Nick
		var file *multipart.FileHeader
		// Get first file from form field "document":
		file, err = c.FormFile("document")
		fn := nick + ":" + strconv.Itoa(int(time.Now().UnixMilli())) + ".ogg"
		// Check for errors:
		if err == nil {
			// 👷 Save file to /uploads directory:
			c.SaveFile(file, fmt.Sprintf("./uploads/%s/%s", room, fn))
			Send <- &models.Msg{Type: "voice", To: room, Msg: fmt.Sprintf("http://127.0.0.1:9393/uploads/%s/%s", room, fn)}
		}

		return
	})
	InitWS()
}

func RunServer(host string) {
	log.Info("Running server on ", host)
	if host == "" {
		host = ":9393"
	}
	log.Error("Server shut down", "error", app.Listen(host))
}
