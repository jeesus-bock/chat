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
		u := new(models.User)
		u.UserID = k
		u.Nick = v.Nick
		// ContainsUser only matches the nick, so we can possibly bail out early
		u.Rooms = v.Rooms
		ret = append(ret, *u)
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
		if utils.ContainsStr(v.Rooms, room) {
			if !utils.ContainsStr(ret, v.Nick) {
				ret = append(ret, v.Nick)
			}
		}
	}
	return
}
