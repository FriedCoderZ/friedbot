package main

import (
	"fmt"

	"github.com/FriedCoderZ/friedbot"
	"github.com/FriedCoderZ/friedbot/plugins/onebot"
	"github.com/FriedCoderZ/friedbot/plugins/test"
)

func main() {
	bot := friedbot.NewBot()
	plugins := []friedbot.Plugin{
		&onebot.API{},
		&test.Plugin{},
	}
	fmt.Print("Installing Plugins...")
	bot.Use(plugins...)
	fmt.Print("OK\n")
	bot.Run(9000)
	fmt.Println("friedbot is up and listening on port 9000")
}
