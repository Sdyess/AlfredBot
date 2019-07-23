package commands

import (
	"strings"

	"time"

	"github.com/bwmarrin/discordgo"
)

//ExecuteCommand Parses and executes the command from the calling code
func ExecuteCommand(s *discordgo.Session, m *discordgo.Message, t0 time.Time) {

	msg := strings.Split(strings.TrimSpace(m.Content), "!")[1]

	if len(msg) > 2 {
		msg = strings.Split(strings.Split(m.Content, " ")[0], "!")[1]
	}

	switch msg {
	case "info":
		HandleInfoCommand(s, m, t0)
	case "ping":
		HandlePingCommand(s, m)
	case "play":
		if !strings.EqualFold(m.Author.Username, "taft") {
			HandleWrongPermissions(s, m, msg)
			return
		}

		game := strings.Join(strings.Split(m.Content, " ")[1:], " ")
		HandlePlayCommand(s, game)
	case "reload":
		HandleReloadCommand(s, m)
	case "purge":

	case "poll":
		pollSplit := strings.Join(strings.Split(m.Content, " ")[1:], " ")
		HandlePollCommand(s, m, pollSplit)
	case "strawpoll":
		titleBegin := strings.Index(m.Content, "{")
		titleEnd := strings.Index(m.Content, "}")
		if titleBegin == -1 || titleEnd == -1 {
			return
		}
		pollTitle := string(m.Content[titleBegin+1:titleEnd])

		strpos := 0
		pollOptionsData := m.Content[titleEnd+1:]
		var pollOptions []string
		for strpos < len(pollOptionsData) {
			test := pollOptionsData[strpos:]
			beginPos := strings.Index(test, "[")
			endPos := strings.Index(test, "]")
			if beginPos == -1 || endPos == -1 || len(test[beginPos:endPos]) <= 0{
				break
			}
			pollOptions = append(pollOptions, test[beginPos+1:endPos])
			strpos = strpos + (endPos + 1)
		}

		if len(pollOptions) < 2 {
			return
		}
		HandleStrawPollCommand(s, m, pollTitle, pollOptions)
	default:
		HandleUnknownCommand(s, m, msg)
	}
}



//HandleUnknownCommand is the default for any commands not listed
func HandleUnknownCommand(s *discordgo.Session, m *discordgo.Message, msg string) {

	c, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		println("Unable to open User Channel: ", err)
		return
	}
	s.ChannelMessageSend(c.ID, "The command ` "+msg+" ` is not recognized.")
}

func HandleWrongPermissions(s *discordgo.Session, m *discordgo.Message, msg string) {

	c, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		println("Unable to open User Channel: ", err)
		return
	}
	s.ChannelMessageSend(c.ID, "The command ` "+msg+" ` is not available to you.")
}


