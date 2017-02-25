package database

import (
	"database/sql"
	"fmt"
)

func LoadDatabaseTimers(db *sql.DB, m *map[int]string) (bool, error) {

	rows, err := db.Query("SELECT * FROM timer_words")
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
	rows, err := db.Query("SELECT * FROM users")
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
	rows, err := db.Query("SELECT * FROM censor_words")
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
