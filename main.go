package main

import (
	"fmt"
	"strings"

	"github.com/AlfredBot/commands"
	"github.com/bwmarrin/discordgo"
)

var (
	Token = "YOUR_TOKEN_HERE"
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

	if m.Content[0] != '!' {
		return
	}

	commands.ExecuteCommand(s, m.Message)

	if m.Content == "!info" {

		message := "```txt\n%s\n%s\n%-16s%-20s\n%-16s%-20s```"
		message = fmt.Sprintf(message, "Debug Information", strings.Repeat("-", len("Debug Information")), "ChannelID", m.ChannelID, "Author", m.Author.Username)
		_, _ = s.ChannelMessageSend(m.ChannelID, message)
	}

	if m.Content == "!deleteme" {
		s.ChannelMessageDelete(m.ChannelID, m.Message.ID)
	}
}
