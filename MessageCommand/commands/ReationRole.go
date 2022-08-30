package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

type ReactionRole struct {
}

func (f ReactionRole) HandleInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionMessageComponent:
		if i.MessageComponentData().CustomID == "select" {
			fmt.Println(i.MessageComponentData().Type())
			fmt.Println(i.MessageComponentData().Values)
			eighteen := "1013963333268942868"
			sixteen := "1013963368329134152"
			twelve := "1013963384745627690"
			six := "1013963402810490951"
			fmt.Println(i.Member.User.Username)

			missing6, missing12, missing16, missing18 := true, true, true, true
			for _, v := range i.MessageComponentData().Values {
				switch v {
				case "6":
					s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, six)
					missing6 = false
				case "12":
					s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, twelve)
					missing12 = false
				case "16":
					s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, sixteen)
					missing16 = false
				case "18":
					s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, eighteen)
					missing18 = false
				}
			}
			//Todo Implement a better System
			if missing6 {
				s.GuildMemberRoleRemove(i.GuildID, i.Member.User.ID, six)
			}
			if missing12 {
				s.GuildMemberRoleRemove(i.GuildID, i.Member.User.ID, twelve)
			}
			if missing16 {
				s.GuildMemberRoleRemove(i.GuildID, i.Member.User.ID, sixteen)
			}
			if missing18 {
				s.GuildMemberRoleRemove(i.GuildID, i.Member.User.ID, eighteen)
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					TTS:             false,
					Content:         "You selected: " + strings.Join(i.MessageComponentData().Values, ", "),
					Components:      nil,
					Embeds:          nil,
					AllowedMentions: nil,
					Files:           nil,
					Flags:           0,
					Choices:         nil,
					CustomID:        "",
					Title:           "",
				},
			})
		}
	}
}

func (f ReactionRole) ReactOnMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.Bot {
		return
	}

	if message.Content == f.Name() {
		//session.ChannelMessageSend(message.ChannelID, "bar")
		session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
			Content: "Hi",
			Embeds:  nil,
			TTS:     false,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							CustomID:    "select",
							Placeholder: "Choose your Roles",
							MinValues:   nil,
							MaxValues:   3,
							Options: []discordgo.SelectMenuOption{
								{
									Label:       "18+",
									Value:       "18",
									Default:     false,
									Description: "18",
								},
								{
									Label:       "16",
									Value:       "16",
									Default:     false,
									Description: "16",
								},
								{
									Label:       "12",
									Value:       "12",
									Default:     false,
									Description: "12",
								},
								{
									Label:       "6",
									Value:       "6",
									Default:     false,
									Description: "6",
								},
							},
							Disabled: false,
						},
					},
				},
			},
			Files:           nil,
			AllowedMentions: nil,
			Reference:       nil,
			File:            nil,
			Embed:           nil,
		})

		return
	}
}

func (f ReactionRole) Name() string {
	return "reactionrole"
}
