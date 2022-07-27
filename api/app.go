package api

import (
	"chat/logger"
	"chat/models"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
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

	// Debug settings allow all origins, TODO change on deployments
	app.Use(cors.New(cors.Config{AllowOrigins: "*"}))

	// The main end-point with all the relevant data
	app.Get("/server", getServerDataHandler(cfg))

	// Temporary debug end-point or maybe use to view data on front.
	app.Get("/rooms", GetRoomsMapHandler)

	// Statically serve uploads dir
	app.Static("/uploads", "./uploads")
	app.Post("/upload/:id/:room", func(c *fiber.Ctx) (err error) {
		log.Info("/upload POST handler")
		room := c.Params("room")
		id := c.Params("id")
		log.Info("Room:", room, "ID: ", id)
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
			_, err = os.Stat("./uploads")
			log.Error(err)
			if errors.Is(err, os.ErrNotExist) {
				log.Info("Creating upload dir ./uploads")
				err = os.Mkdir("./uploads", 0777)
				if err != nil {
					log.Error(err)
				}
				_, err = os.Stat(fmt.Sprintf("./uploads/%s", room))
				if errors.Is(err, os.ErrNotExist) {
					log.Infof("Creating upload room dir %s", fmt.Sprintf("./uploads/%s", room))
					err = os.Mkdir(fmt.Sprintf("./uploads/%s", room), 0777)
					if err != nil {
						log.Error(err)
					}
				}
			}
			err = c.SaveFile(file, fmt.Sprintf("./uploads/%s/%s", room, fn))
			if err != nil {
				c.SendStatus(fiber.StatusBadRequest)
				c.JSON(err)
			}
			Send <- &models.Msg{Type: "voice", To: room, Msg: fmt.Sprintf(cfg.URL+"/uploads/%s/%s", room, fn)}
		} else {
			log.Error(err)
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
