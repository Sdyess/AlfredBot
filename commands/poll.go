package commands

import (
	"github.com/AlfredBot/logger"
	"github.com/bwmarrin/discordgo"
)
// !poll
func HandlePollCommand(s * discordgo.Session, m * discordgo.Message, pollTopic string) {
	message, e := s.ChannelMessageSend(m.ChannelID, pollTopic)
	if e != nil {
		logger.WriteError("A problem occurred while sending a message.", e)
		return
	}

	go s.MessageReactionAdd(m.ChannelID, message.ID, "👎")
	go s.MessageReactionAdd(m.ChannelID, message.ID, "🤷")
	go s.MessageReactionAdd(m.ChannelID, message.ID, "👍")

}

func HandleStrawPollCommand() {

}
