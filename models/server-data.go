package models

type ServerData struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	VoiceURL  string `json:"voiceUrl"`
	UserCount int    `json:"userCount"`
	Rooms     []Room `json:"rooms"`
	Users     []User `json:"users"`
}