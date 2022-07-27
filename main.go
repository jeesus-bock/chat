package main

import (
	"chat/api"
	"chat/logger"
	"chat/models"
	"flag"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func main() {
	log = logger.GetLogger()
	config := new(models.Config)
	host := flag.String("host", "127.0.0.1:9393", "hostname:port")
	name := flag.String("name", "server1", "unique name for server")
	typ := flag.String("type", "MASTER", "Not used atm")
	flag.Parse()
	config.Host = *host
	config.Name = *name
	config.Type = *typ
	log.Infof("Config: %+v", config)
	api.Init(config)
	go pipeline()
	api.RunServer(config.Host)
}
