package models

// This struct can be extended to allow all sorts of information, statuses and descriptions
// and whatnot.
type User struct {
	UserID string   `json:"userId"`
	Nick   string   `json:"nick"`
	Server string   `json:"server"`
	Rooms  []string `json:"rooms"`
}
