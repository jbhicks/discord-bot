package db

import (
	"database/sql"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() error {
	var err error
	DB, err = sql.Open("sqlite3", "./playlists.db")
	if err != nil {
		return err
	}
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS playlists (
		user_id TEXT,
		name TEXT,
		songs TEXT,
		PRIMARY KEY (user_id, name)
	)`)
	return err
}

func CreatePlaylist(userID, name string) error {
	_, err := DB.Exec("INSERT INTO playlists (user_id, name, songs) VALUES (?, ?, ?)", userID, name, "")
	return err
}

func AddToPlaylist(userID, name, url string) error {
	var songs string
	err := DB.QueryRow("SELECT songs FROM playlists WHERE user_id = ? AND name = ?", userID, name).Scan(&songs)
	if err != nil {
		return err
	}
	songs += url + ";"
	_, err = DB.Exec("UPDATE playlists SET songs = ? WHERE user_id = ? AND name = ?", songs, userID, name)
	return err
}

func GetPlaylist(userID, name string) ([]string, error) {
	var songs string
	err := DB.QueryRow("SELECT songs FROM playlists WHERE user_id = ? AND name = ?", userID, name).Scan(&songs)
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimSuffix(songs, ";"), ";"), nil
}

func ListPlaylists(userID string) ([]string, error) {
	rows, err := DB.Query("SELECT name FROM playlists WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var names []string
	for rows.Next() {
		var name string
		rows.Scan(&name)
		names = append(names, name)
	}
	return names, nil
}
