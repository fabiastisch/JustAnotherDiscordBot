package ApplicationCommand

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

type SlashCommandHandler struct {
	cmdMap  map[string]SlashCommand
	GuildID string
	session *discordgo.Session
}

func NewSlashCommandHandler(session *discordgo.Session, guildID string) (handler *SlashCommandHandler) {
	handler = &SlashCommandHandler{
		cmdMap:  make(map[string]SlashCommand),
		GuildID: guildID,
		session: session,
	}
	handler.session.AddHandler(handler.HandleInteractionCreate)
	return
}

func (receiver *SlashCommandHandler) HandleInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if command, ok := receiver.cmdMap[i.ApplicationCommandData().Name]; ok {
		command.Execute(s, i)
	}
}

func (receiver *SlashCommandHandler) RegisterCommand(command SlashCommand) {
	if receiver.cmdMap[command.ApplicationCommand().Name] != nil {
		log.Panicf("Cannot create '%v' command. There is an already existing command.", command.ApplicationCommand().Name)
		return
	}
	receiver.cmdMap[command.ApplicationCommand().Name] = command

	_, err := receiver.session.ApplicationCommandCreate(receiver.session.State.User.ID, receiver.GuildID, command.ApplicationCommand())
	if err != nil {
		log.Panicf("Cannot create '%v' command: %v", command.ApplicationCommand().Name, err)
		return
	}
	log.Printf("Successfully created '%v' command\n", command.ApplicationCommand().Name)
}

func (receiver *SlashCommandHandler) CleanupCommands() {
	for _, v := range receiver.cmdMap {
		err := receiver.session.ApplicationCommandDelete(receiver.session.State.User.ID, receiver.GuildID, v.ApplicationCommand().ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.ApplicationCommand().Name, err)
		}
	}
	log.Println("\nFinished Removing Commands")
}
