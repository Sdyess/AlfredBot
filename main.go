package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/AlfredBot/automod"
	"github.com/AlfredBot/commands"
	"github.com/AlfredBot/database"
	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
)

var (
	Token = ""
	BotID string
)

var db *sql.DB
var err error
var t0 time.Time
var userMap = make(map[uint64]string)

func main() {

	t0 = time.Now()

	db, err = sql.Open("mysql", "username:password@/database")

	if err != nil {
		fmt.Println("Error connecting to database: ", err)
		return
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
	dg.AddHandler(OnGuildMemberAdd)

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection: ", err)
		return
	}

	println("Loading Users...")
	if ok, err := database.LoadDatabaseUsers(db, &userMap); !ok {
		fmt.Println("[ERROR] Issue while loading users table", err)
		return
	}

	println("Running!")

	<-make(chan struct{})
	return
}

//func GuildMemberUpdate()
func OnGuildMemberAdd(s *discordgo.Session, g *discordgo.GuildMemberAdd) {
	if s == nil || g == nil {
		return
	}

	var user = g.User
	if user.ID == BotID {
		return
	}

	st, err := s.UserChannelCreate(user.ID)
	if err != nil {
		return
	}

	//this will be moved to a seperate database loaded messaging system
	s.ChannelMessageSend(st.ID, "Greetings, my name is Alfred. I'm here to help you get adjusted to the chatroom.")
	s.ChannelMessageSend(st.ID, "Before chatting with us, I would appreciate it if you took a moment to review the <#278647380679852032> channel")
	s.ChannelMessageSend(st.ID, "After reviewing the rules, please take another moment and adjust your nickname to your first name.")
	s.ChannelMessageSend(st.ID, "AlfredBot is a discord bot created specifically for this chat and is an open-source project. If you are interested in helping with the creation and advancement of AlfredBot, please visit https://github.com/Sdyessdev/AlfredBot/")

}

func OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if s == nil || m == nil {
		return
	}

	if m.Author.ID == BotID {
		return
	}

	if m.Content == "" {
		go automod.CleanupNudity(s, m.Message)
		return
	}

	if m.Content[0] == '!' && strings.Count(m.Content, "!") < 2 {

		commands.ExecuteCommand(s, m.Message, t0)
		return
	}

	go automod.CleanupNudity(s, m.Message)

	if automod.IsWordCensored(m.Message, db) {
		err := s.ChannelMessageDelete(m.ChannelID, m.Message.ID)
		if err != nil {
			fmt.Println("[Error] Issue deleting a censored message: ", err)
		}
		return
	}

	if automod.IsWordOnTimer(m.Message, db) {
		timer := time.NewTimer(time.Minute)
		go func() {
			<-timer.C
			err := s.ChannelMessageDelete(m.ChannelID, m.ID)
			if err != nil {
				fmt.Println("[Error] Issue deleting a timed message: ", err)
				return
			}
		}()
		return
	}

}
