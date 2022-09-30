package cronStuff

import (
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
	"justAnotherDiscordBot/MessageCommand/commands"
	"log"
	"time"
)

func MensaMenu(session *discordgo.Session) {
	log.Println()
	c := cron.New(cron.WithLocation(time.UTC), cron.WithSeconds())
	//test := "lol"
	/*c.Schedule(cron.Every(time.Second), cron.FuncJob(func() {
		log.Println(test)
	}))*/
	// second minute hour ...+
	mensaChannel := "1025441316592693319"
	RemoveAllChannelMessages(session, mensaChannel)

	imgFile := commands.GetCanteenImageReader()
	defer imgFile.Close()
	session.ChannelMessageSendComplex(mensaChannel, &discordgo.MessageSend{
		Content:    "",
		Embeds:     nil,
		TTS:        false,
		Components: nil,
		Files: []*discordgo.File{
			{
				Name:        "welcome.png",
				ContentType: "image/png",
				Reader:      imgFile,
			},
		},
		AllowedMentions: nil,
		Reference:       nil,
		File:            nil,
		Embed:           nil,
	})

	_, err := c.AddFunc("1 * * * * *", func() {
		log.Println("-")
	})
	if err != nil {
		log.Fatalln(err)
	}
	c.Start()

}

func RemoveAllChannelMessages(session *discordgo.Session, channel string) {

	var messagesIds [][]string
	var messagesIdsTemp []string

	messages, err := session.ChannelMessages(channel, 100, "", "", "")
	if err != nil {
		return
	}
	log.Printf("ChannelMessages: %d\n", len(messages)) // TODO 0 Messages

	for i, m := range messages {
		if (i+1)%100 == 0 {
			//log.Printf("New Array on i: %d\n", i)
			messagesIds = append(messagesIds, messagesIdsTemp)
			messagesIdsTemp = nil
			continue
		}
		messagesIdsTemp = append(messagesIdsTemp, m.ID)
		log.Printf(m.ID)
	}
	messagesIds = append(messagesIds, messagesIdsTemp)

	for _, messages := range messagesIds {
		log.Printf("Remove %d messages", len(messages))
		err := session.ChannelMessagesBulkDelete(channel, messages)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
