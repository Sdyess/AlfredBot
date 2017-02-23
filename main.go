package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/AlfredBot/automod"
	"github.com/AlfredBot/commands"
	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
)

var (
	Token = "YOUR_TOKEN"
	BotID string
)

var db *sql.DB
var err error

func main() {

	db, err = sql.Open("mysql", "username:pass@/db")

	if err != nil {
		fmt.Println("Error connecting to database: ", err)
	}
	fmt.Println("[INFO] Connected to database")

	defer db.Close()

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		println("Error creating discord session: ", err)
		return
	}
	fmt.Println("[INFO] Session Created")

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
		err := s.ChannelMessageDelete(m.ChannelID, m.Message.ID)
		if err != nil {
			fmt.Println("[Error] Issue deleting a censored message: ", err)
		}
	}

	if automod.IsWordOnTimer(m.Message, db) {
		timer := time.NewTimer(time.Minute)
		go func() {
			<-timer.C
			err := s.ChannelMessageDelete(m.ChannelID, m.ID)
			if err != nil {
				fmt.Println("[Error] Issue deleting a timed message: ", err)
			}
		}()
	}

}
