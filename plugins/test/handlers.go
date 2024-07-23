package test

import (
	"log/slog"

	"github.com/FriedCoderZ/friedbot"
)

func eventsLog(ctx *friedbot.Context) {
	events := ctx.GetEvents()
	if events != nil {
		slog.Info("EventsLog", "events", events)
	}
	ctx.Next()
}

func repeat(ctx *friedbot.Context) {
	bot := ctx.GetBot()
	event := ctx.GetEvents().Top()
	group := event.GetGroup()
	msg := event.GetContent()
	_, err := bot.API.SendGroupMsg(group.ID, msg)
	if err != nil {
		slog.ErrorContext(ctx, "repeat", "err", err)
	}
	ctx.Abort()
}
