package commands

import "github.com/bwmarrin/discordgo"

type Foo struct {
}

func (f Foo) ReactOnMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.Bot {
		return
	}

	if message.Content == f.Name() {
		session.ChannelMessageSend(message.ChannelID, "bar")
		return
	}
}

func (f Foo) Name() string {
	return "foo"
}
