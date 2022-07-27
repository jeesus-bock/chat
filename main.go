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

	host := flag.String("host", "127.0.0.1:9393", "hostname:port used as the server")
	url := flag.String("url", "http://127.0.0.1:9393", "hostname:port, used as published url")
	name := flag.String("name", "server1", "unique name for server")
	typ := flag.String("type", "MASTER", "Not used atm")
	noRedis := flag.Bool("noredis", false, "Turn off redis store")
	flag.Parse()

	config.Host = *host
	config.Name = *name
	config.Type = *typ
	config.URL = *url
	config.NoRedis = *noRedis

	log.Infof("Config: %+v", config)
	api.Init(config)
	go pipeline()
	api.RunServer(config.Host)
}
