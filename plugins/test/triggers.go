package test

import (
	"github.com/FriedCoderZ/friedbot"
)

func msgTrigger(ctx *friedbot.Context) bool {
	event := ctx.GetEvents().Top()
	return event.IsMsg()
}

func groupTrigger(ctx *friedbot.Context) bool {
	event := ctx.GetEvents().Top()
	if !event.IsMsg() {
		return false
	}
	msg := event.GetMsg()
	return msg.GroupID == 879423094
}
