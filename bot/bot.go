package bot

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/FriedCoderZ/friedbot/cache"
	"github.com/FriedCoderZ/friedbot/models"
)

type Bot struct {
	API          API
	streams      []Stream
	CacheService *cache.Service
	plugins      []Plugin
}

func (b *Bot) AddStream(stream Stream) {
	b.streams = append(b.streams, stream)
}

func (b *Bot) ClearStreams() {
	b.streams = []Stream{}
}

func (b *Bot) Use(plugins ...Plugin) {
	for _, plugin := range b.plugins {
		err := plugin.Install(b)
		if err != nil {
			slog.Error("failed to Install plugin", "error", err)
			panic("failed to Install plugin")
		}
		b.plugins = append(b.plugins, plugin)
	}
}

func (b *Bot) Run(port int) {
	if b.API == nil {
		panic("no API plugins are used")
	}
	cache.NewService().Start()
	models.NewService("./models/sqlite.db").Start()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
		event := models.NewEvent(requestBody)
		eventCache, err := b.API.ParseEvent(event)
		if err != nil {
			slog.Error("failed to get event cache", "error", err)
			return
		}
		b.handle(eventCache)
	})

	slog.Info("Listening on port", port)
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		slog.Error("ListenAndServe: ", "error", err)
	}
}

func (b *Bot) handle(cache *cache.Cache) {
	for _, stream := range b.streams {
		done, err := stream.Handle(cache)
		if err != nil {
			slog.Error("failed to handle stream", "error", err)
			continue
		}
		if done {
			break
		}
	}
}
