package MessageCommand

import "github.com/bwmarrin/discordgo"

type MessageCommand interface {
	ReactOnMessage(session *discordgo.Session, message *discordgo.MessageCreate)
	Name() string
}

type HandleInteractionCreate interface {
	HandleInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate)
}
