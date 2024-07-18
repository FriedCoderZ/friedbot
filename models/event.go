package models

import (
	"time"
)

type Event struct {
	Data     map[string]any
	CreateAt time.Time
}

func NewEvent(data map[string]any) *Event {
	return &Event{
		Data:     data,
		CreateAt: time.Now(),
	}
}
