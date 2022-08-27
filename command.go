package main

import "github.com/bwmarrin/discordgo"

type Command interface {
	execute(session *discordgo.Session, message *discordgo.MessageCreate)
}
