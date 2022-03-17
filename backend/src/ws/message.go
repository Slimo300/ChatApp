package ws

import "time"

type Message struct {
	Group   uint64    `json:"group"`
	Member  uint64    `json:"member"`
	Message string    `json:"text"`
	Nick    string    `json:"nick"`
	When    time.Time `json:"created"`
}
