package main

import (
	"github.com/FriedCoderZ/friedbot"
	"github.com/FriedCoderZ/friedbot/plugins/onebot"
	"github.com/FriedCoderZ/friedbot/plugins/test"
)

func main() {
	bot := friedbot.NewBot()
	api := onebot.NewAPI(9001)
	plugins := []friedbot.Plugin{
		api,
		&test.Plugin{},
	}
	bot.Use(plugins...)
	bot.Run(9000)
}
