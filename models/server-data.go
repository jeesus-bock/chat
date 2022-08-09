package models

type ServerData struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	URL      string `json:"url"`
	VoiceURL string `json:"voiceUrl"`
}
