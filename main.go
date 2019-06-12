package main

import (
	"database/sql"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/AlfredBot/automod"
	"github.com/AlfredBot/commands"
	"github.com/AlfredBot/database"
	"github.com/AlfredBot/logger"
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
	db = database.Connect()
	if db == nil {
		return
	}

	println("Loading Users...")
	if ok, err := database.LoadDatabaseUsers(db, &userMap); !ok {
		logger.WriteError("Issue while loading users table.", err)
		return
	}

	loaded := automod.LoadAutomodTables(db)
	if !loaded {
		logger.WriteInfo("Automod failed to load tables.")
	}

	defer db.Close()

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		logger.WriteError("Error creating discord session.", err)
		return
	}
	logger.WriteInfo("Session Created.")

	u, err := dg.User("@me")
	if err != nil {
		logger.WriteError("A problem occurred while obtaining account details.", err)
		return
	}

	BotID = u.ID

	dg.AddHandler(messageCreate)
	dg.AddHandler(messageReactionAdd)

	err = dg.Open()
	if err != nil {
		logger.WriteError("A problem occurred while opening a connection.", err)
		return
	}

	logger.WriteInfo("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

//func GuildMemberUpdate()
// func OnGuildMemberAdd(s *discordgo.Session, g *discordgo.GuildMemberAdd) {
// 	if s == nil || g == nil {
// 		return
// 	}

// 	var user = g.User
// 	if user.ID == BotID {
// 		return
// 	}

// 	st, err := s.UserChannelCreate(user.ID)
// 	if err != nil {
// 		return
// 	}
// }

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if s == nil || m == nil {
		return
	}

	if m.Author.ID == BotID {
		return
	}

	if m.Content == "" {
		return
	}

	if m.Content[0] == '!' && strings.Count(m.Content, "!") < 2 {

		commands.ExecuteCommand(s, m.Message, t0)
		return
	}

	if automod.IsWordCensored(m.Message, db) {
		err := s.ChannelMessageDelete(m.ChannelID, m.Message.ID)
		if err != nil {
			logger.WriteError("Issue deleting a censored message.", err)
		}
		return
	}

	if automod.IsWordOnTimer(m.Message, db) {
		timer := time.NewTimer(time.Minute)
		go func() {
			<-timer.C
			err := s.ChannelMessageDelete(m.ChannelID, m.ID)
			if err != nil {
				logger.WriteError("Issue deleting a timed message.", err)
				return
			}
		}()
		return
	}

}

func messageReactionAdd(s * discordgo.Session, reactMsg * discordgo.MessageReactionAdd) {
	_, err := s.ChannelMessage(reactMsg.ChannelID, reactMsg.MessageID)
	if err != nil {
		logger.WriteError("A problem occurred while getting a message.", err)
	}
}
