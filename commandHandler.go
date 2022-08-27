package main

import (
	"github.com/bwmarrin/discordgo"
)

type CommandHandler struct {
	commandFactory *CommandFactory
}

func (commandHandler *CommandHandler) handleCommands(session *discordgo.Session, message *discordgo.MessageCreate) {

	command := commandHandler.commandFactory.getCommand(message.Content)

	if command != nil {
		command.execute(session, message)
	}
}
