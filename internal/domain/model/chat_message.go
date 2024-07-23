package model

import "time"

type Message struct {
	From      string
	To        string
	Content   string
	ChatType  string // private or group
	Timestamp time.Time
}
