package main

import "github.com/bwmarrin/discordgo"

type PingCommand struct{}

func (command *PingCommand) execute(session *discordgo.Session, message *discordgo.MessageCreate) {
	session.ChannelMessageSend(message.ChannelID, "Pong!")
}
