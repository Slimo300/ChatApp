package ws

import "time"

type Message struct {
	Group   uint64
	Member  uint64
	Message string
	When    time.Time
}
