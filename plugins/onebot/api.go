package onebot

import (
	"fmt"

	"github.com/FriedCoderZ/friedbot/cache"
	"github.com/FriedCoderZ/friedbot/core"
	"github.com/FriedCoderZ/friedbot/models"
)

type API struct{}

func (a API) Install(bot *core.Bot) error {
	bot.API = a
	return nil
}

func (API) SendPrivateMsg(userID int64, message string) (messageID int64, err error) {
	//TODO implement me
	panic("implement me")
}

func (API) SendGroupMsg(groupID int64, message string) (messageID int64, err error) {
	//TODO implement me
	panic("implement me")
}

func (API) DeleteMsg(messageID int32) (err error) {
	//TODO implement me
	panic("implement me")
}

func (API) GetMsg(messageID int32) (message any, err error) {
	//TODO implement me
	panic("implement me")
}

func (API) ParseEvent(event *models.Event) (cache *cache.Cache, err error) {
	fmt.Println(event)
	return nil, nil
}
