package main

import (
	"fmt"

	"github.com/FriedCoderZ/friedbot"
	"github.com/FriedCoderZ/friedbot/plugins/chatgpt"
	"github.com/FriedCoderZ/friedbot/plugins/onebot"
)

func main() {
	bot := friedbot.NewBot()
	plugins := []friedbot.Plugin{
		onebot.NewAPI(9001),
		chatgpt.NewPlugin(),
		//&test.Plugin{},
	}
	bot.Use(plugins...)
	fmt.Println("Bot started")
	bot.Run(9000)
}
