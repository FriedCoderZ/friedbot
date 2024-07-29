package chatgpt

import (
	"strings"

	"github.com/FriedCoderZ/friedbot"
)

func msgTrigger(ctx *friedbot.Context) bool {
	event := ctx.GetEvents().Top()
	return event.IsMsg()
}

func prefixTrigger(ctx *friedbot.Context) bool {
	event := ctx.GetEvents().Top()
	msg := event.GetMsg()
	return strings.HasPrefix(msg.Content, triggerPrefix)
}

func groupTrigger(ctx *friedbot.Context) bool {
	event := ctx.GetEvents().Top()
	if !event.IsMsg() {
		return false
	}
	msg := event.GetMsg()
	for _, groupID := range groupWhitelist {
		if groupID == msg.GroupID {
			return true
		}
	}
	return false
}

func slice(ctx *friedbot.Context) bool {
	event := ctx.GetEvents().Top()
	msg := event.GetMsg()
	content := strings.TrimPrefix(msg.Content, triggerPrefix)
	words := strings.Fields(content)
	if len(words) == 0 {
		return false
	}
	ctx.Set("words", words)
	return true
}
