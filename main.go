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

var (
	handlers []*ApplicationCommand.SlashCommandHandler
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

	// Ready Function
	bot.AddHandler(func(s *discordgo.Session, ready *discordgo.Ready) {
		s.UpdateGameStatus(0, "Playing with commands.")
		fmt.Printf("Logged in as: %v#%v\n", s.State.User.Username, s.State.User.Discriminator)
	})

	bot.Identify.Intents = discordgo.IntentGuildMessages
	// connection will receive only events defined by this intent
	// Todo: Add intents if needed

	error = bot.Open()

	if error != nil {
		log.Fatalf("Error while connecting to gateaway: %s", error)
		return
	}
	defer bot.Close()

	for _, guild := range bot.State.Guilds {
		fmt.Println("Create Slash CommandHandler for Guild: " + guild.ID)
		/*x, err := bot.Guild(guild.ID)
		if err != nil {
			return
		}
		fmt.Println(x.Name)*/
		h := ApplicationCommand.NewSlashCommandHandler(bot, guild.ID)
		h.RegisterCommand(commands.Ping{})
		h.RegisterCommand(commands.Canteen{})
		handlers = append(handlers, h)

		log.Print(handlers)
	}

	fmt.Printf("The Bot is now running\n")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	for _, v := range handlers {
		v.CleanupCommands()
	}

	log.Println("Gracefully shutting down.")
}
