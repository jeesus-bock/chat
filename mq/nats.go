package mq

import (
	"chat/logger"
	"chat/models"
	"context"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

var Recv chan *models.Msg
var Send chan *models.Msg
var log *logrus.Logger

func init() {
	log = logger.GetLogger()
	log.Info("Initializing NATS Recv and Send channels")
	Send = make(chan *models.Msg)
	Recv = make(chan *models.Msg)
	nc, err := nats.Connect("nats://foo:bar@127.0.0.1:4222")
	if err != nil {
		log.Fatal("Failed to connect to NATS", "error", err)
	}
	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal("Failed to connect to NATS", "error", err)
	}

	ec.Subscribe("chat", func(msg *models.Msg) {
		Recv <- msg
		log.Info("Received mq msg %+v\n", msg)
	})

	// A go routine for sending messages pushed to the Send channel
	go func() {
		log.Info("Running NATS send goroutine")
		ctx := context.Background()
	outer:
		for {
			select {
			case msg := <-Send:
				log.Info("mqSend", msg)
				ec.Publish("chat", msg)
			case <-ctx.Done():
				break outer
			}
		}
	}()
	log.Info("Initialized NATS")
}

func SendMsg(msg *models.Msg) {
	log.Info("Sending msg to mq")
	Send <- msg
	log.Info("Msg to mq sent")

}
