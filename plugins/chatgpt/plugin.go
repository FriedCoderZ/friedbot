package chatgpt

import (
	"errors"
	"time"

	"github.com/FriedCoderZ/friedbot"
)

type Plugin struct {
	initialized bool
}

func NewPlugin() *Plugin {
	//go checkValid(time.Minute * 30)
	links = make(map[int64]*link)
	return &Plugin{initialized: true}
}

func (p Plugin) Install(bot *friedbot.Bot) error {
	if !p.initialized {
		return errors.New("chatgpt plugin is not init")
	}
	stream := friedbot.NewStream(friedbot.TriggerModeAll)
	stream.AppendTriggers(msgTrigger, prefixTrigger, slice)
	stream.AppendHandlers(tryCreate, tryWithdraw, reply)
	bot.AppendStreams(*stream)
	return nil
}

func checkValid(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		for id, l := range links {
			if !l.isValid() {
				delete(links, id)
			}
		}
	}
	defer ticker.Stop()
}
