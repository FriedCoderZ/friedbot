package models

import (
	"time"
)

type Message struct {
	ID       int64
	Type     string
	SubType  string
	Content  string
	TargetID int64
	Time     time.Time
	Sender   *User
}
