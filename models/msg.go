package models

type Msg struct {
	Type string `json:"type"`
	From string `json:"from"`
	To   string `json:"to"`
	Msg  string `json:"msg"`
	TS   int    `json:"ts"`
}
