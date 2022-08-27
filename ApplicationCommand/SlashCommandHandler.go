package ApplicationCommand

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
)

var (
	cmdMap             map[string]SlashCommand
	registeredCommands []*discordgo.ApplicationCommand
	GuildID            = "1012517914387173437"
)

func RegisterCommand(command SlashCommand) {
	if len(cmdMap) == 0 {
		fmt.Println("Cmd len 0")
		cmdMap = make(map[string]SlashCommand)
	}
	cmdMap[command.ApplicationCommand().Name] = command
}

// FinishCommands
// Bot has to be Running to get `s.State.User.ID`
func FinishCommands(s *discordgo.Session) {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if command, ok := cmdMap[i.ApplicationCommandData().Name]; ok {
			command.Execute(s, i)
		}
	})

	log.Println("Adding commands...")
	registeredCommands = make([]*discordgo.ApplicationCommand, len(cmdMap))

	i := 0
	for _, v := range cmdMap {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, GuildID, v.ApplicationCommand())
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.ApplicationCommand().Name, err)
		}
		registeredCommands[i] = cmd
		i++
	}
}

func RemoveCommands(s *discordgo.Session) {
	for _, v := range registeredCommands {
		err := s.ApplicationCommandDelete(s.State.User.ID, GuildID, v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
	log.Println("\nFinished Removing Commands")
}
