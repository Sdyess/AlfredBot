package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

//Parses and executes the command from the calling code
func ExecuteCommand(s *discordgo.Session, m *discordgo.Message) {

	msg := strings.Split(strings.TrimSpace(m.Content), "!")[1]

	switch msg {
	case "test":
		c, err := s.UserChannelCreate(m.Author.ID)
		if err != nil {
			println("Unable to open User Channel: ", err)
			return
		}
		s.ChannelMessageSend(c.ID, "testm8")
		s.ChannelMessageSend(m.ChannelID, "Hello "+m.Author.Username+" ("+m.ChannelID+")"+m.Author.Email)
	}
}
