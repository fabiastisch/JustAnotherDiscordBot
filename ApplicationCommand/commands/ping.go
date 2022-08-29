package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

type Ping struct {
}

func (e Ping) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Ping Pong :P",
	}
}

func (e Ping) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong! " + fmt.Sprintf("%dms", s.HeartbeatLatency().Milliseconds()),
		},
	})
}
