package MessageCommand

import (
	"github.com/bwmarrin/discordgo"
	"justAnotherDiscordBot/MessageCommand/commands"
	"log"
)

type Handler struct {
	cmdMap map[string]MessageCommand
	//registeredCommands []*discordgo.ApplicationCommand
	//GuildID            string
	session *discordgo.Session
}

func NewHandler(session *discordgo.Session) (handler *Handler) {
	handler = &Handler{
		cmdMap: make(map[string]MessageCommand),
		//GuildID: guildID,
		session: session,
	}
	handler.session.AddHandler(handler.HandleMessageCreate)
	handler.session.AddHandler(handler.HandleInteractionCreate)
	handler.RegisterCommand(commands.ReactionRole{})
	handler.RegisterCommand(commands.Foo{})
	return
}

func (receiver *Handler) HandleInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionMessageComponent:
		for _, v := range receiver.cmdMap {
			if h, ok := v.(HandleInteractionCreate); ok {
				h.HandleInteractionCreate(s, i)
			}
		}
	}
}

func (receiver *Handler) HandleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if command, ok := receiver.cmdMap[m.Content]; ok {
		command.ReactOnMessage(s, m)
	}
}

func (receiver *Handler) RegisterCommand(command MessageCommand) {
	if receiver.cmdMap[command.Name()] != nil {
		log.Panicf("Cannot create '%v' command. There is an already existing command.", command.Name())
		return
	}
	//receiver.registeredCommands = append(receiver.registeredCommands, command)
	receiver.cmdMap[command.Name()] = command
	log.Printf("Successfully created '%v' command\n", command.Name())
}
