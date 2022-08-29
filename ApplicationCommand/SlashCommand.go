package ApplicationCommand

import "github.com/bwmarrin/discordgo"

type SlashCommand interface {
	Execute(s *discordgo.Session, i *discordgo.InteractionCreate)
	ApplicationCommand() *discordgo.ApplicationCommand
}
