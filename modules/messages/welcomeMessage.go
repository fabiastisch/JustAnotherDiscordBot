package messages

import (
	"github.com/bwmarrin/discordgo"
	"image/color"
	"image/png"
	"justAnotherDiscordBot/Picture"
	"log"
	"os"
)

func WelcomeMessage(s *discordgo.Session, e *discordgo.GuildMemberAdd) {
	guild, err := s.Guild(e.GuildID)
	if err != nil {
		log.Panic(err)
	}
	img := Picture.New(400, 200)

	img.Background(color.RGBA{B: 100, A: 255})
	// y := bottomline
	img.AddLabelCenterHorizontal("Willkommen "+e.Member.User.Username, 0, color.White, 28)
	//avatar, err := s.UserAvatar(e.Member.User.ID)
	//img.AddLabel(10, 200, e.Author.AvatarURL("256"))
	//avatar, err := s.UserAvatarDecode(e.Author)
	if err != nil {
		log.Panic(err)
	}

	img.DrawImageBottomCenter(Picture.GetImageFromURL(e.Member.User.AvatarURL("128")))

	//reader := img.ToReader()

	reader, writer, err := os.Pipe()
	defer reader.Close()

	go func() {
		// close the writer, so the reader knows there's no more data
		defer writer.Close()
		if err != nil {
			log.Panic(err)
		}
		//err = png.Encode(writer, img.GetImage())
		err = png.Encode(writer, img.GetImage())

		if err != nil {
			writer.Close()
			log.Panicln(err)
		}
	}()

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
}
