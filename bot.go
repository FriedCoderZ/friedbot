package friedbot

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

type Bot struct {
	API          API
	streams      []Stream
	CacheService *Cache
	plugins      []Plugin
}

func NewBot() *Bot {
	return &Bot{
		API:          nil,
		streams:      []Stream{},
		plugins:      []Plugin{},
		CacheService: NewService(),
	}
}

func (b *Bot) AppendStreams(stream ...Stream) {
	b.streams = append(b.streams, stream...)
}

func (b *Bot) ClearStreams() {
	b.streams = []Stream{}
}

func (b *Bot) Use(plugins ...Plugin) {
	for _, plugin := range plugins {
		err := plugin.Install(b)
		if err != nil {
			slog.Error("failed to Install plugin", "error", err)
			panic("failed to Install plugin")
		}
		b.plugins = append(plugins, plugin)
	}
}

func (b *Bot) Run(port int) {
	if b.API == nil {
		panic("no API plugins are used")
	}
	b.CacheService.Start()
	//models.NewDatabase("./models/sqlite.db").Start()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				slog.Error("failed to close request body", "error", err)
			}
		}(r.Body)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		requestBody := make(map[string]any)
		if err := json.Unmarshal(body, &requestBody); err != nil {
			http.Error(w, "Failed to parse request body", http.StatusBadRequest)
			return
		}
		event, err := b.API.ParseEvent(requestBody)
		if err != nil {
			slog.Error("failed to parse event", "error", err)
			return
		}
		cache, err := b.API.GetCache(event)
		if err != nil {
			slog.Error("failed to get event cache", "error", err)
			return
		}
		cache.Add(event)
		b.handle(cache)
	})

	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		slog.Error("ListenAndServe: ", "error", err)
	}
}

func (b *Bot) handle(cache *Ring) {
	for _, stream := range b.streams {
		ctx := NewContext(b, &stream, cache)
		triggered := stream.Handle(ctx)
		if triggered {
			break
		}
	}
}

func (b *Bot) Reply(msg *Message, content string) (msgID int64, err error) {
	if msg.Type == MsgTypeGroup && msg.GroupID != 0 {
		msgID, err = b.API.SendGroupMsg(msg.GroupID, content)
	} else if msg.Type == MsgTypePrivate && msg.User != nil && msg.User.ID != 0 {
		msgID, err = b.API.SendPrivateMsg(msg.User.ID, content)
	} else {
		return 0, fmt.Errorf("message does not have enough data to reply")
	}
	return msgID, err
}
