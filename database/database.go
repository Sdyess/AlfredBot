package database

import (
	"database/sql"
	"fmt"
)

type BotInfo struct {
	status int
	last   string
}

const loadDatabaseTimersSQL = "SELECT * FROM TIMER_WORDS"
const loadDatabaseUsersSQL = "SELECT * FROM USERS"
const loadDatabaseCensorSQL = "SELECT * FROM CENSOR_WORDS"
const loadDatabaseBotInfoSQL = "SELECT * FROM BOT_INFO"

func Connect() *sql.DB {
	db, err := sql.Open("mysql", "root:root@/alfredbot")
	if err != nil {
		fmt.Println("[ERROR] Unable to connect to database: ", err)
		return nil
	}

	fmt.Println("[INFO] Connected to database.")
	return db
}

func LoadDatabaseTimers(db *sql.DB, m *map[int]string) (bool, error) {
	fmt.Println("[INFO] Loading Removable Words...")
	rows, err := db.Query(loadDatabaseTimersSQL)
	if err != nil {
		return false, err
	}

	defer rows.Close()

	for rows.Next() {
		var id int
		var word string
		err = rows.Scan(&id, &word)
		if err != nil {
			return false, err
		}

		(*m)[id] = word
	}

	fmt.Println("[INFO] Removable Words loaded.")
	return true, nil
}

func LoadDatabaseUsers(db *sql.DB, m *map[uint64]string) (bool, error) {
	rows, err := db.Query(loadDatabaseUsersSQL)
	if err != nil {
		return false, err
	}

	defer rows.Close()

	for rows.Next() {
		var id uint64
		var word string
		err = rows.Scan(&id, &word)
		if err != nil {
			return false, err
		}

		(*m)[id] = word
	}

	fmt.Println("[INFO] Users loaded.")
	return true, nil
}

func LoadDatabaseCensoredWords(db *sql.DB, m *map[int]string) (bool, error) {
	rows, err := db.Query(loadDatabaseCensorSQL)
	if err != nil {
		return false, err
	}

	defer rows.Close()

	for rows.Next() {
		var id int
		var word string
		err = rows.Scan(&id, &word)
		if err != nil {
			return false, err
		}

		(*m)[id] = word
	}

	fmt.Println("[INFO] Censored Words loaded.")
	return true, nil
}

func LoadBotInfo(db *sql.DB) (bool, BotInfo, error) {

	var info BotInfo

	rows, err := db.Query(loadDatabaseBotInfoSQL)
	if err != nil {
		return false, info, err
	}

	defer rows.Close()

	for rows.Next() {
		var lastPlaying string
		var status int
		err = rows.Scan(&lastPlaying, &status)

		if err != nil {
			return false, info, err
		}

		info.status = status
		info.last = lastPlaying
	}

	fmt.Println("[INFO] Bot Info loaded.")
	return true, info, nil

}
