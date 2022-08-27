package main

import "github.com/bwmarrin/discordgo"

type PatrickCommand struct{}

func (command *PatrickCommand) execute(session *discordgo.Session, message *discordgo.MessageCreate) {
	session.ChannelMessageSend(message.ChannelID, "Nein, hier ist Patrick!")
}
