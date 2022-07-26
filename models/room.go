package models

type Room struct {
	Topic string   `json:"topic"`
	Users []string `json:"users"`
}
