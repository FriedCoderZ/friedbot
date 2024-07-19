package main

import (
	"github.com/FriedCoderZ/friedbot/core"
	"github.com/FriedCoderZ/friedbot/plugins/onebot"
)

func main() {
	bot := core.NewBot()
	plugins := []core.Plugin{
		&onebot.API{},
	}
	bot.Use(plugins...)
	bot.Run(9000)
}
