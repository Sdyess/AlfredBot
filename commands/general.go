package commands

import (
	"fmt"
	"github.com/AlfredBot/automod"
	"github.com/bwmarrin/discordgo"
	"strings"
	"time"
)

//HandleInfoCommand is the !info command
func HandleInfoCommand(s *discordgo.Session, m *discordgo.Message, t0 time.Time) {

	t1 := time.Now()
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		fmt.Println("[ERROR] Issue finding channel by ID: ", err)
		return
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
		return
	}
}

func HandleReloadCommand(s *discordgo.Session, m *discordgo.Message) {
	automod.ReloadTables()
	s.ChannelMessageSend(m.ChannelID, "Database tables have been reloaded!")
}

func HandlePurgeCommand() {

}
