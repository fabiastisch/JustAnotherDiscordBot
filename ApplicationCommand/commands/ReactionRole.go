package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"reflect"
	"sort"
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

	applicationCommandOptions = append(applicationCommandOptions, &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionBoolean,
		Name:        "role-option-removable",
		Description: "Removable",
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

	removable := false
	if opt, ok := optionMap["role-option-removable"]; ok {
		removable = opt.BoolValue()
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
			}
		}
		maxValues := 1
		if opt, ok := optionMap["role-option-maxvalue"]; ok {
			maxValues = int(opt.IntValue())
		}

		integerOptionMinValue := 0

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
								CustomID: func() string {
									if removable {
										return "ApplicationReactionRole-select-removable"
									} else {
										return "ApplicationReactionRole-select"
									}
								}(),
								Placeholder: "Choose your Roles",
								MinValues:   &integerOptionMinValue,
								MaxValues:   maxValues,
								Options:     selectMenuOptions, // max 25 (https://discord.com/developers/docs/interactions/message-components#select-menu-object-select-menu-structure)
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

func (e ReactionRole) HandleInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionMessageComponent:
		if i.MessageComponentData().CustomID != "ApplicationReactionRole-select" && i.MessageComponentData().CustomID != "ApplicationReactionRole-select-removable" {
			return
		}
		removable := i.MessageComponentData().CustomID == "ApplicationReactionRole-select-removable"
		allRoles := make(map[string]string)
		for _, component := range i.Message.Components {
			if component.Type() == discordgo.ActionsRowComponent {
				actionRow := component.(*discordgo.ActionsRow)
				//log.Println("Actionrow: " + fmt.Sprint(actionRow))
				for _, messageComponent := range actionRow.Components {
					if messageComponent.Type() == discordgo.SelectMenuComponent {
						//log.Printf("Component: %#v\n", messageComponent.(*discordgo.SelectMenu))
						menu := messageComponent.(*discordgo.SelectMenu)
						//log.Println("Options: " + fmt.Sprint(menu.Options))
						for _, option := range menu.Options {
							allRoles[option.Value] = option.Label
						}
					}
				}
			}
		}
		log.Println(reflect.TypeOf(e).Name() + " | InteractionMessageComponent: " + i.MessageComponentData().CustomID + " | User: " + i.Member.User.Username)

		// get all GuildMember roles and sort it.
		guildMemberRoles := i.Member.Roles
		sort.Strings(guildMemberRoles)

		var roles []string
		// add selected Roles
		for _, v := range i.MessageComponentData().Values {
			roles = append(roles, v)
			if UserHasRole(guildMemberRoles, v) {
				continue
			}
			log.Println("Add GuildMemberRole " + v + " | " + i.Member.User.Username)
			err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, v)
			if err != nil {
				log.Println(err)
			}
			//missings = false
		}

		//make a map with the selected roles and rolename
		gRoles, err := s.GuildRoles(i.GuildID)
		if err != nil {
			panic(err)
		}
		selecetRoles := make(map[string]string)
		for _, r := range roles {
			for _, v := range gRoles {
				if r == v.ID {
					selecetRoles[v.ID] = v.Name
					break
				}
			}
		}

		responseMessage := ""
		for _, roleName := range selecetRoles {
			responseMessage += "+ " + roleName + "\n"
		}

		removes := map[string]string{}
		if removable {
			removes = RemoveUnsetRoles(s, i, &allRoles, &selecetRoles, guildMemberRoles)
			for _, roleName := range removes {
				responseMessage += "- " + roleName + "\n"
			}
		}

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			// Type Value must be one of {4, 5, 6, 7, 9}
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Components: nil,
				Flags:      discordgo.MessageFlagsEphemeral,
				Embeds: []*discordgo.MessageEmbed{
					{
						URL:         "",
						Type:        "rich",
						Title:       "Updated Roles",
						Description: "```diff\n" + responseMessage + "```",
						Timestamp:   "",
						Color:       0,
						Footer:      nil,
						Image:       nil,
						Thumbnail:   nil,
						Video:       nil,
						Provider:    nil,
						Author:      nil,
						Fields:      []*discordgo.MessageEmbedField{},
					},
				},
			},
		})
		if err != nil {
			panic(err)
		}

	}
}

func UserHasRole(userRoles []string, role string) bool {
	i := sort.SearchStrings(userRoles, role)
	return i < len(userRoles) && userRoles[i] == role
}

// RemoveUnsetRoles Remove the roles that are not in selectedRoles but in allRoles.
//
// Returns a Map with roleid and roleName with removed roles.
func RemoveUnsetRoles(s *discordgo.Session, i *discordgo.InteractionCreate, allRoles *map[string]string, selectedRoles *map[string]string, guildMemberRoles []string) map[string]string {
	response := make(map[string]string)
	for key, value := range *allRoles {
		// Check if all Role is in selectedRole, if not remove role
		if _, ok := (*selectedRoles)[key]; !ok {
			// Check if the User has the role, if not continue
			if !UserHasRole(guildMemberRoles, key) {
				continue
			}
			err := s.GuildMemberRoleRemove(i.GuildID, i.Member.User.ID, key)
			if err != nil {
				log.Panicf("Error with GuildMemberRoleRemove:\n%v", err)
			}
			response[key] = value
		}
	}
	return response
}
