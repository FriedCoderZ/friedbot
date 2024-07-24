package chatgpt

import (
	"log/slog"

	"github.com/FriedCoderZ/friedbot"
)

func messageLog(ctx *friedbot.Context) {
	events := ctx.GetEvents()
	event := ctx.GetEvents().Top()
	msg := event.GetMsg()
	if events != nil {
		slog.Info("msg log", "msg", msg)
	}
	ctx.Next()
}

func repeat(ctx *friedbot.Context) {
	bot := ctx.GetBot()
	event := ctx.GetEvents().Top()
	msg := event.GetMsg()
	if msg == nil {
		slog.ErrorContext(ctx, "not message event", "event", event)
		ctx.Abort()
		return
	}
	_, err := bot.API.SendGroupMsg(msg.GroupID, msg.Content)
	if err != nil {
		slog.ErrorContext(ctx, "repeat", "err", err)
	}
	ctx.Abort()
}
