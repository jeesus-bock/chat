package api

import (
	"chat/models"
	"chat/utils"
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

var Recv chan *models.Msg
var Send chan *models.Msg

type WSConn struct {
	Conn *websocket.Conn
	Room string
	Nick string
}
type Container struct {
	mu        sync.Mutex
	sendConns map[string]WSConn
}

var topics map[string]string
var cc *Container

func init() {
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

func (c *Container) sendMsg(id string, msg *models.Msg) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.sendConns[id].Conn.WriteJSON(msg)
}

// GetUsers returns a map of all users on the server.
// TODO this needs defining and adjustments.
func GetUsers(cfg models.Config, conns map[string]WSConn) (ret map[string]*models.User) {
	ret = make(map[string]*models.User)
	rooms := make([]string, 0)
	for k, v := range conns {
		ret[k] = new(models.User)
		ret[k].UserID = k
		ret[k].Nick = v.Nick
		// This is left borken, we need to define ws conn <-> room relation
		if !utils.ContainsStr(rooms, v.Room) {
			rooms = append(rooms, v.Room)
		}
		ret[k].Rooms = rooms
		ret[k].Server = cfg.Host
	}
	return
}
