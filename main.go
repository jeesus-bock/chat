package main

import (
	"chat/logger"
	"chat/models"
	"chat/mq"
	"chat/ws"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func main() {
	log = logger.GetLogger()
	go func() {
		// Create ticker to send ping, debug code.
		ticker := time.NewTicker(2099 * time.Second)
		for {
			select {
			case mqMsg := <-mq.Recv:
				log.Debug("Received NATS message: ", fmt.Sprintf("%+v\n", mqMsg))
				if mqMsg.Type == "nick" {
					msg := mqMsg.From + " is now known as " + mqMsg.Msg
					mqMsg.From = "server"
					mqMsg.Msg = msg
				}
				ws.Send <- mqMsg
			case wsMsg := <-ws.Recv:
				log.Debug("Received WS message", fmt.Sprintf("%+v\n", wsMsg))
				mq.Send <- wsMsg
			case t := <-ticker.C:
				log.Debugf("Sending keepalive", "t", t)
				mq.Send <- &models.Msg{Type: "ping", TS: int(time.Now().UnixMilli()), To: "*", Msg: "Keepalive"}
			}
		}
	}()
	ws.RunServer()
}
