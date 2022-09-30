package commands

import (
	"github.com/bwmarrin/discordgo"
	"justAnotherDiscordBot/clients/canteenClient"
)

type Canteen struct {
}

func (e Canteen) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "mensa",
		Description: "Erhalte Informationen zu den kommenden Mahlzeiten",
	}
}

func (e Canteen) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	canteenClient.Request(s, i)
}
