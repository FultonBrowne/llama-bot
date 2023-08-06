package main

import(
	"DiscordBot/bot"
	"DiscordBot/config"
	"fmt"
)

func main() {
	err := config.ReadConfig()

	if err != nil {
		fmt.Println(err.Error())
	}

	bot.Start()

	<-make(chan struct{})
	return
}