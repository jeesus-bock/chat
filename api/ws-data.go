package api

import (
	"chat/models"
	"chat/utils"
)

// GetUsers returns a map of all users on the server.
// TODO this needs defining and adjustments.
func GetUsers(cfg *models.Config, conns map[string]WSConn) (ret []models.User) {
	ret = make([]models.User, 0)
	for k, v := range conns {
		tmp := new(models.User)
		tmp.UserID = k
		tmp.Nick = v.Nick
		// ContainsUser only matches the nick, so we can possibly bail out early
		if !utils.ContainsUser(ret, *tmp) {
			// This is left borken, we need to define ws conn <-> room relation
			rooms := make([]string, 0)
			if !utils.ContainsStr(rooms, v.Room) {
				rooms = append(rooms, v.Room)
			}
			tmp.Rooms = rooms
			tmp.Server = cfg.Host
			ret = append(ret, *tmp)
		}
	}
	return
}

func GetRooms(cfg *models.Config, conns map[string]WSConn) (ret []models.Room) {
	log.Info("GetRooms", rooms)
	ret = make([]models.Room, 0)
	for k, v := range rooms {
		rr := new(models.Room)
		rr.Name = k
		rr.Topic = v
		rr.Users = GetRoomUsers(k, conns)
		// Only insert each room once
		if !utils.ContainsRoom(ret, *rr) {
			ret = append(ret, *rr)
		}
	}
	return
}
func GetRoomsMap() map[string]string {
	return rooms
}

func GetRoomUsers(room string, conns map[string]WSConn) (ret []string) {
	ret = make([]string, 0)
	for _, v := range conns {
		if v.Room == room {
			if !utils.ContainsStr(ret, v.Nick) {
				ret = append(ret, v.Nick)
			}
		}
	}
	return
}
