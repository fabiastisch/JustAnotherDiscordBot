package timed

import (
	"github.com/bwmarrin/discordgo"
	"time"
)

type Timer struct {
	session *discordgo.Session
}

func New(session *discordgo.Session) *Timer {
	t := &Timer{
		session: session,
	}
	return t
}

func (t *Timer) Schedule(call func(s *discordgo.Session, data ...any), interval time.Duration, data any) *time.Ticker {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			if data != nil {
				call(t.session, data)
			} else {
				call(t.session)
			}
		}
	}()
	return ticker
}
