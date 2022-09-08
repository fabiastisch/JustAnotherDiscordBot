package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"reflect"
	"strings"
)

type ReactionRole struct {
}

func (e ReactionRole) ApplicationCommand() *discordgo.ApplicationCommand {
	integerOptionMinValue := 1.0
	var applicationCommandOptions []*discordgo.ApplicationCommandOption

	applicationCommandOptions = append(applicationCommandOptions, &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionRole,
		Name:        "role-option",
		Description: "Role option",
		Required:    true,
	})

	for i := 1; i < 10; i++ {
		applicationCommandOptions = append(applicationCommandOptions, &discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionRole,
			Name:        "role-option" + fmt.Sprint(i),
			Description: "Role option",
			Required:    false,
		})
	}
	applicationCommandOptions = append(applicationCommandOptions, &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionInteger,
		Name:        "role-option-maxvalue", // lowercase only!
		Description: "Maximal auswählbar",
		MinValue:    &integerOptionMinValue,
		MaxValue:    11,
		Required:    false,
	})

	return &discordgo.ApplicationCommand{
		Name:        "reactionrole",
		Description: "create a new Reaction role...",
		Options:     applicationCommandOptions,
	}
}

func (e ReactionRole) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {

	/*err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: nil,
	})
	if err != nil {
		log.Fatalln(err)
	}*/
	// Access options in the order provided by the user.
	options := i.ApplicationCommandData().Options

	// Or convert the slice into a map
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	// make a set to ignore dupplicates in roles
	allRoles := make(map[string]bool)

	if opt, ok := optionMap["role-option"]; ok {
		var selectMenuOptions []discordgo.SelectMenuOption

		role := opt.RoleValue(s, i.GuildID)
		allRoles[role.ID] = true // add to set
		selectMenuOptions = append(selectMenuOptions, discordgo.SelectMenuOption{
			Label:       role.Name,
			Value:       role.ID,
			Default:     false,
			Description: "-",
		})

		for j := 1; j < 10; j++ {
			if opt, ok := optionMap["role-option"+fmt.Sprint(j)]; ok {
				role := opt.RoleValue(s, i.GuildID)
				if _, ok := allRoles[role.ID]; ok { // check if role is in set
					continue
				}
				allRoles[role.ID] = true // add to set

				selectMenuOptions = append(selectMenuOptions, discordgo.SelectMenuOption{
					Label:       role.Name,
					Value:       role.ID,
					Default:     false,
					Description: "-",
				})
			} else {
				break
			}

		}
		maxValues := 1
		if opt, ok := optionMap["role-option-maxvalue"]; ok {
			maxValues = int(opt.IntValue())
		}

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			//Type Value must be one of {9, 4, 5}
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				TTS:     false,
				Content: "",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.SelectMenu{
								CustomID:    "ApplicationReactionRole-select",
								Placeholder: "Choose your Roles",
								MinValues:   nil,
								MaxValues:   maxValues,
								Options:     selectMenuOptions,
								Disabled:    false,
							},
						},
					},
				},
				Embeds:          nil,
				AllowedMentions: nil,
				Files:           nil,
				Flags:           0,
				Choices:         nil,
				CustomID:        "",
				Title:           "",
			},
		})
		if err != nil {
			log.Panicln(err)
		}

	}
}

func (f ReactionRole) HandleInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionMessageComponent:
		if i.MessageComponentData().CustomID != "ApplicationReactionRole-select" {
			return
		}
		log.Println(reflect.TypeOf(f).Name() + " | InteractionMessageComponent: " + i.MessageComponentData().CustomID + " | User: " + i.Member.User.Username)

		fmt.Println(i.MessageComponentData().Type())
		fmt.Println(i.MessageComponentData().Values)
		log.Println(i.Data)

		//fmt.Println(i.Member.User.Username)

		//s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, i.MessageComponentData())

		//missings := make([]bool, len(i.MessageComponentData().Values))

		var roles []string

		for _, v := range i.MessageComponentData().Values {
			log.Println("Add GuildMemberRole " + v + " | " + i.Member.Nick)
			err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, v)
			if err != nil {
				log.Println(err)
			}
			roles = append(roles, v)
			//missings = false

		}

		gRoles, err := s.GuildRoles(i.GuildID)
		if err != nil {
			panic(err)
		}
		var selecetRoles []string
		for _, r := range roles {
			for _, v := range gRoles {
				if r == v.ID {
					selecetRoles = append(selecetRoles, v.Name)
					break
				}
			}
		}
		//Todo Implement a better System
		/*if missing {
			s.GuildMemberRoleRemove(i.GuildID, i.Member.User.ID, six)
		}*/

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			// Type Value must be one of {4, 5, 6, 7, 9}
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				TTS:             false,
				Content:         "Du hast die Rolle/n erhalten: " + strings.Join(selecetRoles, ","),
				Components:      nil,
				Embeds:          nil,
				AllowedMentions: nil,
				Files:           nil,
				Flags:           discordgo.MessageFlagsEphemeral,
				Choices:         nil,
				CustomID:        "",
				Title:           "",
			},
		})
		if err != nil {
			panic(err)
		}

		/*s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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
		})*/
	}
}
