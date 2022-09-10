package Interfaces

import "github.com/bwmarrin/discordgo"

type HandleInteractionCreate interface {
	HandleInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate)
}
