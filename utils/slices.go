package utils

import "chat/models"

func ContainsStr(ss []string, s string) bool {
	for _, str := range ss {
		if s == str {
			return true
		}
	}
	return false
}
func ContainsRoom(rr []models.Room, r models.Room) bool {
	for _, room := range rr {
		if room.Name == r.Name {
			return true
		}
	}
	return false
}

func ContainsUser(uu []models.User, u models.User) bool {
	for _, user := range uu {
		if user.Nick == u.Nick {
			return true
		}
	}
	return false
}
