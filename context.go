package friedbot

import (
	"log/slog"
	"time"
)

type Context struct {
	index   int
	bot     *Bot
	stream  *Stream
	events  *Ring
	data    map[string]any
	aborted bool
	err     error
}

func NewContext(bot *Bot, stream *Stream, events *Ring) *Context {
	return &Context{
		index:  -1,
		bot:    bot,
		stream: stream,
		events: events,
		data:   make(map[string]any),
	}
}

func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	return
}

func (ctx *Context) Done() <-chan struct{} {
	return nil
}

func (ctx *Context) Err() error {
	return ctx.err
}

func (ctx *Context) Value(key any) any {
	if _, ok := key.(string); !ok {
		return nil
	}
	return ctx.data[key.(string)]
}

func (ctx *Context) Set(key string, value any) {
	ctx.data[key] = value
}

func (ctx *Context) Next() {
	if ctx.aborted {
		return
	}
	if ctx.index >= len(ctx.stream.handlers)-1 {
		if ctx.index != -1 {
			slog.Error("no next handler", "handler", ctx.stream.handlers[ctx.index])
			return
		} else {
			slog.Error("no handler", "stream", ctx.stream)
			return
		}
	}
	ctx.index++
	ctx.stream.handlers[ctx.index](ctx)
}

func (ctx *Context) Abort() {
	ctx.aborted = true
}

func (ctx *Context) GetEvents() *Ring {
	return ctx.events
}

func (ctx *Context) GetBot() *Bot {
	return ctx.bot
}
