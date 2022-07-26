package main

import (
	"chat/api"
	"chat/models"
	"chat/mq"
	"fmt"
	"time"
)

func pipeline() {
	// Create ticker to send ping, debug code.
	ticker := time.NewTicker(2099 * time.Second)
	for {
		select {
		case mqMsg := <-mq.Recv:
			log.Debug("Received NATS message: ", fmt.Sprintf("%+v\n", mqMsg))
			api.Send <- mqMsg
		case wsMsg := <-api.Recv:
			log.Debug("Received WS message", fmt.Sprintf("%+v\n", wsMsg))
			mq.Send <- wsMsg
		case t := <-ticker.C:
			log.Debugf("Sending keepalive", "t", t)
			mq.Send <- &models.Msg{Type: "ping", TS: int(time.Now().UnixMilli()), To: "*", Msg: "Keepalive"}
		}
	}
}
