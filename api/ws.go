package api

import (
	"chat/models"
	"chat/store"
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
	Conn  *websocket.Conn
	Rooms []string
	Nick  string
}
type Container struct {
	mu        sync.Mutex
	sendConns map[string]WSConn
}

// the in-memory store for message history, used if -noredis flag is set
var oldMessages []models.Msg

var rooms map[string]string
var cc *Container

func InitWS(cfg *models.Config) {
	cc = new(Container)
	cc.sendConns = make(map[string]WSConn)
	Recv = make(chan *models.Msg)
	Send = make(chan *models.Msg)
	rooms = make(map[string]string)
	rooms["main"] = "Main chatroom"
	rooms["testing"] = "Testing grounds"
	rooms["test2"] = "Testing grounds #2"

	app.Use("/", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// The ws listening path contains :room parameter for different chatrooms
	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		id := c.Params("id")
		nick := c.Query("nick")

		if nick == "" {
			log.Error("Empty nick")
		}
		cc.sendConns[id] = WSConn{Conn: c, Nick: nick, Rooms: make([]string, 0)}
		log.Infof("New websocket connection, nick: %s id: s", nick, id)
		var (
			mt     int
			rawMsg []byte
			err    error
		)

		// Test using redis backend
		var msgList []models.Msg
		if cfg.NoRedis {
			msgList = oldMessages
		} else {
			msgList, err = store.GetLocalHistory()
			if err != nil {
				log.Error("error loading local msg history from redis: ", err)
			}
		}
		log.Infof("Sending %d old messages", len(msgList))

		cc.sendMsg(id, &models.Msg{Type: "connected", TS: int(time.Now().UnixMilli()), From: "server", To: cfg.Name})

		// Another infinite loop for reading messages and sending then to Recv channel.
		for {
			nick := cc.sendConns[id].Nick
			if mt, rawMsg, err = c.ReadMessage(); err != nil {
				// These aren't actually errors as ws.close() gives us error.
				log.Debug("Error reading websocket message:", "error", err)
				lRooms := cc.sendConns[id].Rooms
				delete(cc.sendConns, id)
				// Send leave msg to all the channels the user was on.
				for _, r := range lRooms {
					Send <- &models.Msg{Type: "leave", TS: int(time.Now().UnixMilli()), To: r, From: "server", Msg: nick}
				}
				break
			}
			msg := new(models.Msg)
			log.Debug("JSON ws message: ", string(rawMsg))
			json.Unmarshal(rawMsg, msg)
			log.Debug("Received websocket message", "data", msg, "type", strconv.Itoa(mt))

			// message to change the user nick
			if msg.Type == "nick" {
				if entry, ok := cc.sendConns[id]; ok {
					entry.Nick = msg.Msg
					cc.sendConns[id] = entry
				}
			}
			if msg.Type == "topic" {
				rooms[msg.To] = msg.Msg
			}
			if msg.Type == "join" {
				if !utils.ContainsStr(cc.sendConns[id].Rooms, msg.To) {
					if entry, ok := cc.sendConns[id]; ok {
						entry.Rooms = append(entry.Rooms, msg.To)
						cc.sendConns[id] = entry
					}
				}
				if _, ok := rooms[msg.To]; !ok {
					rooms[msg.To] = "Default topic"
				}
				cc.sendMsg(id, &models.Msg{Type: "start_history", To: msg.To, TS: int(time.Now().UnixMilli())})
				for _, m := range msgList {
					if m.To == msg.To {
						cc.sendMsg(id, &m)
					}
				}
				cc.sendMsg(id, &models.Msg{Type: "end_history", To: msg.To, TS: int(time.Now().UnixMilli())})
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
			if cfg.NoRedis {
				oldMessages = append(oldMessages, *msg)
			} else {
				store.AddLocalHistoryMsg(*msg)
			}
			// Send messages to this room and to wildcard room "*"
			for k, v := range cc.sendConns {
				if utils.ContainsStr(v.Rooms, msg.To) || msg.To == "*" {
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
