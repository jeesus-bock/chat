package store

import (
	"chat/logger"
	"chat/models"
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

var LOCAL_HISTORY = "LOCAL_HISTORY"
var rdb *redis.Client
var log *logrus.Logger

func init() {
	log = logger.GetLogger()
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	log.Debug("rdb: ", rdb)
}
func AddLocalHistoryMsg(msg models.Msg) (err error) {
	log.Debugf("Adding msg to redis: %+v", msg)
	ctx := context.Background()
	jsonStr, err := json.Marshal(msg)
	if err != nil {
		return
	}
	_, err = rdb.RPush(ctx, LOCAL_HISTORY, jsonStr).Result()
	return
}
func GetLocalHistory() (ret []models.Msg, err error) {
	log.Debug("Getting local msg history from redis")
	ctx := context.Background()
	list, err := rdb.LRange(ctx, LOCAL_HISTORY, 0, -1).Result()
	if err != nil {
		return
	}
	for _, val := range list {
		var um = new(models.Msg)
		err = json.Unmarshal([]byte(val), um)
		if err != nil {
			return
		}
		ret = append(ret, *um)
	}
	log.Debugf("%d messages loaded from redis", len(ret))
	return
}
