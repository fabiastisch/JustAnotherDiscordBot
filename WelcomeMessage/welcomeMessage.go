package WelcomeMessage

import (
	"github.com/bwmarrin/discordgo"
	"image/color"
	"image/png"
	"justAnotherDiscordBot/Picture"
	"log"
	"os"
)

func Handler(s *discordgo.Session, e *discordgo.GuildMemberAdd) {
	/*if e.Author.Bot {
		return
	}*/

	go SendWelcomeMessage(s, e)
}

func SendWelcomeMessage(s *discordgo.Session, e *discordgo.GuildMemberAdd) {
	log.Println("New GuildMemberAdd ")
	guild, err := s.Guild(e.GuildID)
	if err != nil {
		log.Panic(err)
	}
	log.Println("Test")

	img := Picture.Image(400, 200)
	log.Println("Test")

	img.Background(color.RGBA{0, 0, 100, 255})
	log.Println("Test")
	// y := bottomline
	img.AddLabelCenterHorizontal("Willkommen "+e.Member.User.Username, 0, color.White)
	//avatar, err := s.UserAvatar(e.Member.User.ID)
	//img.AddLabel(10, 200, e.Author.AvatarURL("256"))
	//avatar, err := s.UserAvatarDecode(e.Author)
	log.Println("Test")
	if err != nil {
		log.Panic(err)
	}

	img.DrawImageBottomCenter(Picture.GetImageFromURL(e.Member.User.AvatarURL("128")))
	log.Println("Test")

	//reader := img.ToReader()

	reader, writer, err := os.Pipe()

	go func() {
		// close the writer, so the reader knows there's no more data
		defer writer.Close()
		log.Println("ToReader")
		if err != nil {
			log.Panic(err)
		}
		log.Println("Encode...")
		//err = png.Encode(writer, img.GetImage())
		err = png.Encode(writer, img.GetImage())
		log.Println("Encoded")

		if err != nil {
			writer.Close()
			log.Panicln(err)
		}
		log.Println("Encodeasd")
	}()

	log.Println("x")

	log.Println("Test")
	_, err = s.ChannelMessageSendComplex(guild.SystemChannelID, &discordgo.MessageSend{
		Content:    "",
		Embeds:     nil,
		TTS:        false,
		Components: nil,
		Files: []*discordgo.File{
			{
				Name:        "welcome.png",
				ContentType: "image/png",
				Reader:      reader,
			},
		},
		AllowedMentions: nil,
		Reference:       nil,
		File:            nil,
		Embed:           nil,
	})
	if err != nil {
		log.Panic(err)
	}
	defer reader.Close()
	log.Println("F")
}
