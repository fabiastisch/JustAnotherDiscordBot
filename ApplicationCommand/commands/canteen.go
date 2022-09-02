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
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: canteenClient.Request(),
		},
	})
}
