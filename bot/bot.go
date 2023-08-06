package bot

import (
	"DiscordBot/config"
	"DiscordBot/utils"
	"fmt"
	"context"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"github.com/jmorganca/ollama/api"
)
var botID string
var client *discordgo.Session
var apiClient  *api.Client
func Start() {
	session, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println(err)
		return
	}
	session.AddHandler(message)
	session.AddHandler(ready)
	apiClient = api.NewClient()
	fmt.Print("Bot is online")
	defer session.Close()
	if err = session.Open(); err != nil {
		fmt.Println(err)
		return
	}

	scall := make(chan os.Signal, 1)
	signal.Notify(scall, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSEGV, syscall.SIGHUP)
	<-scall
}

func ready(bot *discordgo.Session, event *discordgo.Ready) {
	guildsSize := len(bot.State.Guilds)
	bot.UpdateGameStatus(0, strconv.Itoa(guildsSize) + " with your emotions")
}

func IsDM(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		return false
	}
	if channel.Type != discordgo.ChannelTypeDM {
		return false
	}
	return true
}

func message(bot *discordgo.Session, message *discordgo.MessageCreate) {
	/* Discord bot commands */
	if message.Author.Bot { return }
	
	if strings.HasPrefix(message.Content, "&") || IsDM(bot, message){
		fmt.Println("content:", message.Content)
		ping := bot.HeartbeatLatency().Truncate(60)
		if message.Content == "&ping" {
			bot.ChannelMessageSend(message.ChannelID,`My latency is **` + ping.String() + `**!`)
		} else if message.Content == "&author" {
			bot.ChannelMessageSend(message.ChannelID, "My author is Gonz#0001, I'm only a template discord bot made in golang.")
		} else if message.Content == "&github" {
			embed := embed.NewEmbed().
				SetAuthor(message.Author.Username, message.Author.AvatarURL("1024")).
				SetThumbnail(message.Author.AvatarURL("1024")).
				SetTitle("My repository").
				SetDescription("You can find my repository by clicking [here](https://github.com/gonzyui/Discord-Template).").
				SetColor(0x00ff00).MessageEmbed
			bot.ChannelMessageSendEmbed(message.ChannelID, embed)
		} else if message.Content == "&botinfo" {
			guilds := len(bot.State.Guilds)
			embed := embed.NewEmbed().
				SetTitle("My informations").
				SetDescription("Some informations about me :)").
				AddField("GO version:", runtime.Version()).
				AddField("DiscordGO version:", discordgo.VERSION).
				AddField("Concurrent tasks:", strconv.Itoa(runtime.NumGoroutine())).
				AddField("Latency:", ping.String()).
				AddField("Total guilds:", strconv.Itoa(guilds)).MessageEmbed
			bot.ChannelMessageSendEmbed(message.ChannelID, embed)
		} else {
			var latest api.GenerateResponse
			var output string
			generateContext := []int{} // TODO: Get conversation context
			request := api.GenerateRequest{Model: "llama2", Prompt: message.Content, Context: generateContext}
			fn := func(response api.GenerateResponse) error {

				latest = response
				output = output + response.Response

				fmt.Print(response.Response)
				return nil
			}
			if err := apiClient.Generate(context.Background(), &request, fn); err != nil {
				bot.ChannelMessageSend(message.ChannelID, "There seems to be a problem, please cry to the developer (FultonBrowne on github)")
				return
			}
			bot.ChannelMessageSend(message.ChannelID, output)
			latest.Summary()
		}
	}
	
}
