package chatgpt

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/FriedCoderZ/friedbot"
)

func reply(ctx *friedbot.Context) {
	event := ctx.GetEvents().Top()
	msg := event.GetMsg()
	l := getLinkByMsg(msg)
	content := strings.TrimPrefix(msg.Content, triggerPrefix)
	l.addMsg(&message{
		Content: content,
		Role:    roleTypeUser,
	})
	generateMsg, err := l.generateMsg()
	if err != nil {
		slog.ErrorContext(ctx, "chatgpt generate message error", "err", err)
	}
	if generateMsg != nil {
		if generateMsg.Role != roleTypeError {
			l.addMsg(generateMsg)
		}
		fmt.Println(l.messages)
		bot := ctx.GetBot()
		_, err = bot.Reply(msg, generateMsg.Content)
		if err != nil {
			slog.ErrorContext(ctx, "chatgpt send message error", "err", err)
		}
	}
	ctx.Abort()
}

func tryCreate(ctx *friedbot.Context) {
	word, _ := ctx.Value("words").([]string)
	flag := false
	for _, w := range triggerWordCreate {
		if w == word[0] {
			flag = true
		}
	}
	if !flag {
		ctx.Next()
		return
	} else {
		event := ctx.GetEvents().Top()
		msg := event.GetMsg()
		l := getLinkByMsg(msg)
		m := make([]message, len(basicMessages))
		copy(m, basicMessages)
		l.messages = m
		bot := ctx.GetBot()
		_, err := bot.Reply(msg, "已重建新会话")
		if err != nil {
			slog.ErrorContext(ctx, "chatgpt send message error", "err", err)
		}
		ctx.Abort()
	}
}

func tryWithdraw(ctx *friedbot.Context) {
	word, _ := ctx.Value("words").([]string)
	flag := false
	for _, w := range triggerWordWithdraw {
		if w == word[0] {
			flag = true
		}
	}
	if !flag {
		ctx.Next()
		return
	} else {
		event := ctx.GetEvents().Top()
		msg := event.GetMsg()
		l := getLinkByMsg(msg)
		bot := ctx.GetBot()
		var r string
		if len(l.messages) < 2 {
			r = "已经没有消息可以撤回了"
		} else {
			l.withdrawMsg()
			r = "已回撤上一条对话"
		}
		_, err := bot.Reply(msg, r)
		if err != nil {
			slog.ErrorContext(ctx, "chatgpt send message error", "err", err)
		}
		ctx.Abort()
	}
}
