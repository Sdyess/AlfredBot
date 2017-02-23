package automod

import (
	"strings"

	"fmt"

	"database/sql"

	"github.com/bwmarrin/discordgo"
)

var timerMap = make(map[int]string)

func LoadDatabaseTimers(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM timer_words")

	if err != nil {
		fmt.Println("[ERROR] Issue loading censored words: ", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var word string
		err = rows.Scan(&id, &word)
		if err != nil {
			fmt.Printf("[ERROR] Issue reading rows: ", err)
		}

		timerMap[id] = word
	}

	fmt.Println("[INFO] Censored Words loaded.")
}

func IsWordCensored(m *discordgo.Message) bool {
	//this will check through a preloaded map eventually
	var words [3]string
	words[0] = "dudu"
	words[1] = "brained"
	words[2] = "dorf"

	tokens := strings.Split(m.Content, " ")
	for i := 0; i < len(words); i++ {
		for j := 0; j < len(tokens); j++ {
			if strings.EqualFold(words[i], tokens[j]) {
				return true
			}
		}
	}
	return false
}

func IsWordOnTimer(m *discordgo.Message, db *sql.DB) bool {

	if len(timerMap) == 0 {
		fmt.Println("Loading timer words table...")
		LoadDatabaseTimers(db)
	}

	tokens := strings.Split(m.Content, " ")
	for _, v := range timerMap {
		for j := 0; j < len(tokens); j++ {
			if strings.EqualFold(v, tokens[j]) {
				fmt.Printf("[LOG] Message queued to be erased: %s", m.Content)
				return true
			}
		}
	}
	return false
}
