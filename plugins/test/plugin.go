package test

import (
	"github.com/FriedCoderZ/friedbot"
)

type Plugin struct {
}

func (t Plugin) Install(bot *friedbot.Bot) error {
	stream := friedbot.NewStream(friedbot.TriggerModeAll)
	stream.AppendTriggers(msgTrigger, groupTrigger)
	stream.AppendHandlers(eventsLog, repeat)
	bot.AppendStreams(*stream)
	return nil
}
