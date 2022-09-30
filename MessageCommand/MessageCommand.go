package MessageCommand

import "github.com/bwmarrin/discordgo"

type MessageCommand interface {
	ReactOnMessage(session *discordgo.Session, message *discordgo.MessageCreate)
	Name() string
}
