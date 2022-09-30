package commands

import (
	"github.com/bwmarrin/discordgo"
	"time"
)

type Remember struct {
}

func (e Remember) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "remember",
		Description: "Sends a Reminder message",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "remember-time-string",
				Description: "When to remember (e.g. 1m,1h,1d)",
				Required:    true,
			},
		},
	}
}

func (e Remember) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options

	// Or convert the slice into a map
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	response := "I got You"

	if opt, ok := optionMap["remember-time-string"]; ok {
		optionDuration := opt.StringValue()
		duration, err := time.ParseDuration(optionDuration)
		if err != nil {
			response = "Could not parse the Duration"
		} else {
			time.AfterFunc(duration, func() {
				s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
					Content:         "" + i.Member.User.Mention(),
					Embeds:          nil,
					TTS:             false,
					Components:      nil,
					Files:           nil,
					AllowedMentions: nil,
					Reference:       nil,
					File:            nil,
					Embed:           nil,
				})
			})
		}

	} else {
		response = "There was a problem..."
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})

}
