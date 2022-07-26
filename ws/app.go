package ws

import (
	"chat/logger"
	"chat/models"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
	"github.com/sirupsen/logrus"
)

var Recv chan *models.Msg
var Send chan *models.Msg

var log *logrus.Logger
var app *fiber.App

type WSConn struct {
	Conn *websocket.Conn
	Room string
	Nick string
}
type Container struct {
	mu        sync.Mutex
	sendConns map[string]WSConn
}

func (c *Container) sendMsg(id string, msg *models.Msg) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.sendConns[id].Conn.WriteJSON(msg)
}

var topics map[string]string
var cc *Container

func init() {
	log = logger.GetLogger()
	cc = new(Container)
	cc.sendConns = make(map[string]WSConn)
	topics = make(map[string]string)
	app = fiber.New()
	Recv = make(chan *models.Msg)
	Send = make(chan *models.Msg)
	app.Use(cors.New(cors.Config{AllowOrigins: "*"}))
	app.Get("/monitor", monitor)
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
	app.Use("/", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// The ws listening path contains :room parameter for different chatrooms
	app.Get("/ws/:id/:room", websocket.New(func(c *websocket.Conn) {
		room := c.Params("room")
		id := c.Params("id")
		nick := c.Query("nick")

		roomObj := new(models.Room)
		for _, v := range cc.sendConns {
			if v.Room == room {
				roomObj.Users = append(roomObj.Users, v.Nick)
			}
		}
		roomObj.Topic = topics[room]
		Send <- &models.Msg{Type: "join", To: room, From: "server", TS: int(time.Now().UnixMilli()), Msg: nick}
		cc.sendConns[id] = WSConn{Conn: c, Room: room, Nick: nick}
		log.Info("New websocket connection", "id", id, "room", room)
		var (
			mt     int
			rawMsg []byte
			err    error
		)

		roomJson, err := json.Marshal(roomObj)
		if err != nil {
			log.Error("Failed to marshal roomObj")
		}
		cc.sendMsg(id, &models.Msg{Type: "connected", TS: int(time.Now().UnixMilli()), From: "server", To: room, Msg: string(roomJson)})

		// Another infinite loop for reading messages and sending then to Recv channel
		for {
			nick := cc.sendConns[id].Nick
			if mt, rawMsg, err = c.ReadMessage(); err != nil {
				log.Debug("Error reading websocket message:", "error", err)
				delete(cc.sendConns, id)
				Send <- &models.Msg{Type: "leave", TS: int(time.Now().UnixMilli()), To: room, From: "server", Msg: nick}
				break
			}
			msg := new(models.Msg)
			json.Unmarshal(rawMsg, msg)
			log.Debug("Received websocket message", "data", msg, "type", strconv.Itoa(mt))

			// message to change the user nick
			if msg.Type == "nick" {
				if entry, ok := cc.sendConns[id]; ok {
					entry.Nick = msg.Msg
					cc.sendConns[id] = entry
				}
			}
			// Replace the From field with user's nickname
			msg.From = nick
			// Send the unmarshaled message to the recv channel
			Recv <- msg
		}

	}))
	// A go routine for sending messages pushed to the Send channel
	go func() {
		for {
			msg := <-Send
			// Send messages to this room and to wildcard room "*"
			for k, v := range cc.sendConns {
				if msg.To == v.Room || msg.To == "*" {
					log.Debug("Sending ws message", "id", k, "msg", msg)
					cc.sendMsg(k, msg)
				}
			}
		}
	}()
	log.Info("Initialized websocket listening")
}

func RunServer() {
	log.Info("Running server")
	log.Error("Server shut down", "error", app.Listen(":9393"))
}
