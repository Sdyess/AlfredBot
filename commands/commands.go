/* The plans for this part of the bot are on hold until the MySQL drivers begin to be used, so custom commands can be
 * made a bit more easily as well as having proper permission checks per command
 */

package commands

import (
	"fmt"
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

		game := strings.Split(m.Content, " ")
		var newGame string
		for i, tmpGame := range game {
			if i != 0 {
				newGame += tmpGame + " "
			}
		}

		HandlePlayCommand(s, newGame)
	//case "purge":
	//if(m.Author.)
	default:
		HandleUnknownCommand(s, m, msg)
	}
}

/*func ExecuteCommandWithArgs() {

}*/

/*
func hasPermission() bool{

}
*/

//HandleInfoCommand is the !info command
func HandleInfoCommand(s *discordgo.Session, m *discordgo.Message, t0 time.Time) {

	t1 := time.Now()
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		fmt.Println("[ERROR] Issue finding channel by ID: ", err)
	}

	channelName := channel.Name
	message := "```txt\n%s\n%s\n%-16s%-20s\n%-16s%-20s\n%-16s%-20s```"
	message = fmt.Sprintf(message, "AlfredBot Information", strings.Repeat("-", len("AlfredBot Information")), "ChannelID", m.ChannelID, "Channel Name", channelName, "Uptime", (t1.Sub(t0).String()))
	s.ChannelMessageSend(m.ChannelID, message)
}

//HandlePingCommand is for !ping
func HandlePingCommand(s *discordgo.Session, m *discordgo.Message) {

	s.ChannelMessageSend(m.ChannelID, "pong")
}

//HandlePlayCommand sets now playing status for the bot
func HandlePlayCommand(s *discordgo.Session, game string) {
	err := s.UpdateStatus(0, game)
	if err != nil {
		println("[Error] Issue while updating bot status: ", err)
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
