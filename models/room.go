package models

type Room struct {
	Name  string   `json:"name"`
	Topic string   `json:"topic"`
	Users []string `json:"users"`
}
