package main

import (
	"fmt"
	"strings"

	"time"

	"github.com/AlfredBot/automod"
	"github.com/AlfredBot/commands"
	"github.com/bwmarrin/discordgo"
)

var (
	Token = "MjgyNDU2MTIyMzM2NzM5MzMw.C4wSNA.ySm0QJtH7cxrVK5UFwy2_56_IGs"
	BotID string
)

func main() {
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		println("Error creating discord session: ", err)
		return
	}

	u, err := dg.User("@me")
	if err != nil {
		println("Error obtaining account details: ", err)
		return
	}

	BotID = u.ID

	dg.AddHandler(OnMessageCreate)
	//dg.AddHandler(GuildMemberUpdate)

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection: ", err)
	}

	println("Running!")

	<-make(chan struct{})
	return
}

//func GuildMemberUpdate()
func OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotID {
		return
	}

	if m.Content[0] == '!' && strings.Count(m.Content, "!") < 2 {
		commands.ExecuteCommand(s, m.Message)
	}

	if automod.IsWordCensored(m.Message) {
		s.ChannelMessageSend(m.ChannelID, "NO")
		s.ChannelMessageDelete(m.ChannelID, m.ID)
	}

	if automod.IsWordOnTimer(m.Message) {
		timer := time.NewTimer(time.Minute)
		go func() {
			<-timer.C
			s.ChannelMessageDelete(m.ChannelID, m.ID)
		}()
	}

}
