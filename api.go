package friedbot

import "time"

type API interface {
	SendPrivateMsg(userID int64, msg string) (msgID int64, err error)
	SendGroupMsg(groupID int64, msg string) (msgID int64, err error)
	GetMsg(msgID int32) (msg *Message, err error)
	ParseEvent(requestBody map[string]any) (event Event, err error)
	GetCache(event Event) (cache *Ring, err error)
	// DeleteMsg(msgID int32) (err error)
}

type Event interface {
	GetData() map[string]any
	GetTime() time.Time
	IsMsg() bool
	GetMsg() *Message
	GetContent() string
	GetUser() *User
	GetGroup() *Group
}
