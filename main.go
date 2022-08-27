package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"justAnotherDiscordBot/ApplicationCommand"
	"justAnotherDiscordBot/ApplicationCommand/commands"
	"log"
	"os"
	"os/signal"
	"syscall"
)

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

	bot.AddHandler(reactOnMessage)
	bot.Identify.Intents = discordgo.IntentGuildMessages
	// connection will receive only events defined by this intent
	// Todo: Add intents if needed

	error = bot.Open()

	if error != nil {
		log.Fatalf("Error while connecting to gateaway: %s", error)
		return
	}

	ApplicationCommand.RegisterCommand(commands.Ping{}) // could be done in ApplicationCommand package
	ApplicationCommand.FinishCommands(bot)

	fmt.Printf("The Bot is now running\n")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	bot.Close()

	ApplicationCommand.RemoveCommands(bot)
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
