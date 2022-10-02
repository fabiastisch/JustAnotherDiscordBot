package cronStuff

import (
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
	"justAnotherDiscordBot/MessageCommand/commands"
	"log"
	"time"
)

func MensaMenu(session *discordgo.Session) {

	c := cron.New(cron.WithLocation(time.UTC), cron.WithSeconds())
	//test := "lol"
	/*c.Schedule(cron.Every(time.Second), cron.FuncJob(func() {
		log.Println(test)
	}))*/
	// second minute hour ...+

	_, err := c.AddFunc("0 1 15 * * *", func() { // 15:01 Clock

		log.Println("Call Mensa cron")
		mensaChannel := "1025441316592693319"
		RefreshMensaMenu(session, mensaChannel, 7)
	})
	log.Println("Added Mensa Cron")
	mensaChannel := "1025441316592693319"
	RefreshMensaMenu(session, mensaChannel, 7)
	if err != nil {
		log.Fatalln(err)
	}
	c.Start()

}

func RefreshMensaMenu(session *discordgo.Session, channelID string, daysDays int) {
	RemoveAllChannelMessages(session, channelID)

	today := time.Now()
	if today.Hour() >= 15 {
		today = today.AddDate(0, 0, 1)
	}
	for i := daysDays - 1; i >= 0; i-- {
		imgFile := commands.GetCanteenImageReader(today.AddDate(0, 0, i))
		if imgFile != nil {
			defer imgFile.Close()
			session.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
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
		}
	}
	/*imgFile := commands.GetCanteenImageReader(today.AddDate(0, 0, 2))
	defer imgFile.Close()
	session.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
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
	*/
}

func RemoveAllChannelMessages(session *discordgo.Session, channel string) {

	var messagesIds [][]string
	var messagesIdsTemp []string

	messages, err := session.ChannelMessages(channel, 100, "", "", "")
	if err != nil {
		return
	}
	//log.Printf("ChannelMessages: %d\n", len(messages))

	for i, m := range messages {
		if (i+1)%100 == 0 {
			//log.Printf("New Array on i: %d\n", i)
			messagesIds = append(messagesIds, messagesIdsTemp)
			messagesIdsTemp = nil
			continue
		}
		messagesIdsTemp = append(messagesIdsTemp, m.ID)
		//log.Printf(m.ID)
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
