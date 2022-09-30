package messages

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

func LeaveMessage(s *discordgo.Session, e *discordgo.GuildMemberRemove) {
	guild, err := s.Guild(e.GuildID)
	if err != nil {
		log.Panic(err)
	}

	_, err = s.ChannelMessageSendComplex(guild.SystemChannelID, &discordgo.MessageSend{
		Content:         e.User.Mention() + " hat uns leider verlassen.",
		Embeds:          nil,
		TTS:             false,
		Components:      nil,
		AllowedMentions: nil,
		Reference:       nil,
		File:            nil,
		Embed:           nil,
	})
	if err != nil {
		log.Panic(err)
	}
}
