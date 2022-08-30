package ApplicationCommand

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

type SlashCommandHandler struct {
	cmdMap             map[string]SlashCommand
	registeredCommands []*discordgo.ApplicationCommand
	GuildID            string
	session            *discordgo.Session
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
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		if command, ok := receiver.cmdMap[i.ApplicationCommandData().Name]; ok {
			command.Execute(s, i)
		}
	}
}

func (receiver *SlashCommandHandler) RegisterCommand(command SlashCommand) {
	if receiver.cmdMap[command.ApplicationCommand().Name] != nil {
		log.Panicf("Cannot create '%v' interaction. There is an already existing interaction.", command.ApplicationCommand().Name)
		return
	}

	applicationCommand, err := receiver.session.ApplicationCommandCreate(receiver.session.State.User.ID, receiver.GuildID, command.ApplicationCommand())
	if err != nil {
		log.Panicf("Cannot create '%v' interaction: %v", command.ApplicationCommand().Name, err)
		return
	}
	receiver.registeredCommands = append(receiver.registeredCommands, applicationCommand)
	receiver.cmdMap[command.ApplicationCommand().Name] = command
	log.Printf("Successfully created '%v' interaction\n", command.ApplicationCommand().Name)
}

func (receiver *SlashCommandHandler) CleanupCommands() {
	for _, v := range receiver.registeredCommands {
		err := receiver.session.ApplicationCommandDelete(receiver.session.State.User.ID, receiver.GuildID, v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' interaction: %v", v.Name, err)
		}
	}
	log.Println("\nFinished removing interactions")
}
