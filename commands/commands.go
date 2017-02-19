package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

//ExecuteCommand Parses and executes the command from the calling code
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
		s.ChannelMessageSend(m.ChannelID, "Hello "+m.Author.Username+" ("+m.ChannelID+")")
	case "info":
		message := "```txt\n%s\n%s\n%-16s%-20s\n%-16s%-20s```"
		message = fmt.Sprintf(message, "Debug Information", strings.Repeat("-", len("Debug Information")), "ChannelID", m.ChannelID, "Author", m.Author.Username)
		_, _ = s.ChannelMessageSend(m.ChannelID, message)
	default:
		c, err := s.UserChannelCreate(m.Author.ID)
		if err != nil {
			println("Unable to open User Channel: ", err)
			return
		}
		s.ChannelMessageSend(c.ID, "The command ` "+msg+" ` is not recognized.")
	}
}
