package chatgpt

import (
	"github.com/FriedCoderZ/friedbot"
)

var groupWhitelist = []int64{
	879423094,
}

type Plugin struct {
}

func NewPlugin() *Plugin {
	return &Plugin{}
}

func (t Plugin) Install(bot *friedbot.Bot) error {
	stream := friedbot.NewStream(friedbot.TriggerModeAll)
	stream.AppendTriggers(msgTrigger, groupTrigger)
	stream.AppendHandlers(messageLog, repeat)
	bot.AppendStreams(*stream)
	return nil
}
