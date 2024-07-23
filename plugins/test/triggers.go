package test

import (
	"fmt"

	"github.com/FriedCoderZ/friedbot"
)

func testTrigger(ctx *friedbot.Context) bool {
	fmt.Println("test trigger")
	fmt.Println(ctx)
	return true
}

func msgTrigger(ctx *friedbot.Context) bool {
	event := ctx.GetEvents().Top()
	return event.IsMsg()
}

func groupTrigger(ctx *friedbot.Context) bool {
	event := ctx.GetEvents().Top()
	if !event.IsMsg() {
		return false
	}
	group := event.GetGroup()
	if group == nil {
		return false
	}
	fmt.Println(group.ID)
	if group.ID == 879423094 {
		return true
	}
	return false
}
