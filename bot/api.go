package bot

import (
	"github.com/FriedCoderZ/friedbot/cache"
	"github.com/FriedCoderZ/friedbot/models"
)

type API interface {
	SendPrivateMsg(userID int64, message string) (messageID int64, err error)
	SendGroupMsg(groupID int64, message string) (messageID int64, err error)
	DeleteMsg(messageID int32) (err error)
	GetMsg(messageID int32) (message any, err error)
	ParseEvent(event *models.Event) (cache *cache.Cache, err error)
}
