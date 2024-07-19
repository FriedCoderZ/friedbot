package main

import (
	"fmt"

	"github.com/FriedCoderZ/friedbot/core"
	"github.com/FriedCoderZ/friedbot/plugins/onebot"
)

func main() {
	bot := core.NewBot()
	plugins := []core.Plugin{
		&onebot.API{},
	}
	fmt.Print("Installing Plugins...")
	bot.Use(plugins...)
	fmt.Print("OK\n")
	fmt.Println("Running Bot...")
	bot.Run(9000)
}
