package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"justAnotherDiscordBot/ApplicationCommand"
	"justAnotherDiscordBot/ApplicationCommand/commands"
	"justAnotherDiscordBot/MessageCommand"
	. "justAnotherDiscordBot/modules"
	"log"
	"os"
	"os/signal"
	"reflect"
	"syscall"
)

var (
	handlers []*ApplicationCommand.SlashCommandHandler
)

func init() {
	//Configure log
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	discordToken := os.Getenv("DISCORD_TOKEN")

	bot, error := discordgo.New("Bot " + discordToken)

	if error != nil {
		log.Fatalf("Error while starting bot: %s", error)
	}

	// Ready Function
	bot.AddHandler(func(s *discordgo.Session, ready *discordgo.Ready) {
		s.UpdateGameStatus(0, "Playing with commands.")
		log.Printf("Logged in as: %v#%v\n", s.State.User.Username, s.State.User.Discriminator)
	})

	bot.AddHandler(reactOnMessage)

	//bot.Identify.Intents = discordgo.IntentGuildMessages
	bot.Identify.Intents = discordgo.IntentsAll
	// connection will receive only events defined by this intent
	// Todo: Add intents if needed
	bot.AddHandler(WelcomeMessage)
	bot.AddHandler(LOGHANDLER)
	error = bot.Open()

	if error != nil {
		log.Fatalf("Error while connecting to gateaway: %s", error)
		return
	}
	defer bot.Close()

	for _, guild := range bot.State.Guilds {
		log.Println("Create Slash CommandHandler for Guild: " + guild.ID)
		x, err := bot.Guild(guild.ID)
		if err != nil {
			return
		}
		log.Println(x.Name)
		h := ApplicationCommand.NewSlashCommandHandler(bot, guild.ID)
		h.RegisterCommand(commands.Ping{})
		h.RegisterCommand(commands.ReactionRole{})
		handlers = append(handlers, h)
	}

	MessageCommand.NewHandler(bot)

	log.Printf("The Bot is now running\n")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	for _, v := range handlers {
		v.CleanupCommands()
	}

	log.Println("Gracefully shutting down.")
}

func reactOnMessage(session *discordgo.Session, message *discordgo.MessageCreate) {

	if message.Author.Bot {
		return
	}

	if message.Content == "ping" {
		session.ChannelMessageSend(message.ChannelID, "pong")
		return
	}

	if message.Content == "krossekrabbe" {
		session.ChannelMessageSend(message.ChannelID, "Nein, hier ist Patrick!")
		return
	}

}

func LOGHANDLER(s *discordgo.Session, i interface{}) {
	color := "\u001b[33m"
	reset := "\u001b[0m"
	log.Println(color + "Event: " + fmt.Sprint(reflect.TypeOf(i)) + reset)
}
