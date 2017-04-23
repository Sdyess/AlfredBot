package automod

import (
	"strings"

	"database/sql"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"

	"time"

	"github.com/AlfredBot/database"
	"github.com/bwmarrin/discordgo"
	"github.com/go-nude"
)

var removeableWordsMap = make(map[int]string)
var censoredWordsMap = make(map[int]string)

func ReloadTables() {
	for k := range removeableWordsMap {
		delete(removeableWordsMap, k)
	}

	for k := range censoredWordsMap {
		delete(censoredWordsMap, k)
	}

	fmt.Println("[INFO] Word cache cleared.")
}

//IsWordCensored (* discordgo.Message) bool
//Words that match this check are immediately removed from chat
func IsWordCensored(m *discordgo.Message, db *sql.DB) bool {
	if len(censoredWordsMap) == 0 {
		fmt.Println("Loading censored words table...")
		if ok, _ := database.LoadDatabaseCensoredWords(db, &censoredWordsMap); !ok {
			return false
		}
	}

	tokens := strings.Split(m.Content, " ")
	for i, v := range censoredWordsMap {
		for j := 0; j < len(tokens); j++ {
			if _, ok := censoredWordsMap[i]; !ok {
				fmt.Println("[ERROR] Attempt to access index out of bounds during censor search")
				return false
			}

			if strings.EqualFold(v, tokens[j]) {
				fmt.Printf("\n[LOG] Message erased: %s", m.Content)
				return true
			}
		}
	}
	return false
}

func IsWordOnTimer(m *discordgo.Message, db *sql.DB) bool {

	if len(removeableWordsMap) == 0 {
		fmt.Println("Loading removeable words table...")
		if ok, _ := database.LoadDatabaseTimers(db, &removeableWordsMap); !ok {
			return false
		}
	}

	tokens := strings.Split(m.Content, " ")
	for i, v := range removeableWordsMap {
		for j := 0; j < len(tokens); j++ {
			if _, ok := removeableWordsMap[i]; !ok {
				fmt.Println("[ERROR] Attempt to access index out of bounds during removable search")
				return false
			}

			if strings.EqualFold(v, tokens[j]) {
				fmt.Printf("\n[LOG] Message queued to be erased: %s", m.Content)
				return true
			}
		}
	}
	return false
}

func CleanupNudity(s *discordgo.Session, m *discordgo.Message) {

	var url string

	for _, j := range m.Embeds {
		if j == nil {
			fmt.Println("[ERROR]: ", j)
			return
		}
		url = j.URL
		response, err := http.Get(url)
		if err != nil {
			fmt.Println("[ERROR]: ", err)
			return
		}

		file, err := os.Create(GenerateFileName())
		if err != nil {
			fmt.Println("[ERROR] Unable to create file ", err)
			return
		}

		_, err = io.Copy(file, response.Body)
		if err != nil {
			fmt.Println("[ERROR]: Unable to copy image to file", err)
			return
		}

		if val, err := nude.IsNude(file.Name()); val {
			if err != nil {
				fmt.Println("[ERROR] ", err)
				return
			}
			err := s.ChannelMessageDelete(m.ChannelID, m.ID)
			if err != nil {
				return
			}
		}
		err = response.Body.Close()
		if err != nil {
			return
		}
		err = file.Close()
		if err != nil {
			return
		}
		err = os.Remove(file.Name())
		if err != nil {
			return
		}

	}
}

func GenerateFileName() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())

	b := make([]rune, 5)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
