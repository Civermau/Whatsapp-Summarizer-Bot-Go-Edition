package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Database operations
func initDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "work.db")
	if err != nil {
		return nil, err
	}

	// Create messages table if it doesn't exist
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		chat_id TEXT NOT NULL,
		sender TEXT NOT NULL,
		message TEXT,
		message_type TEXT,
		timestamp DATETIME
	);`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// insertMessage inserts a message into the database
func insertMessage(chatID, sender, message, messageType string, timestamp time.Time) error {
	query := `INSERT INTO messages (chat_id, sender, message, message_type, timestamp) VALUES (?, ?, ?, ?, ?)`
	_, err := db.Exec(query, chatID, sender, message, messageType, timestamp)
	return err
}

func getMessages(chatID string, limit int) (string, error) {
	query := `SELECT sender, message FROM messages WHERE chat_id = ? ORDER BY timestamp DESC LIMIT ?`
	rows, err := db.Query(query, chatID, limit)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var messages []string
	for rows.Next() {
		var sender, message string
		if err := rows.Scan(&sender, &message); err != nil {
			return "", err
		}
		messages = append(messages, fmt.Sprintf("%s: %s", sender, message))
	}
	if err := rows.Err(); err != nil {
		return "", err
	}

	// Reverse the messages to get chronological order (oldest to newest)
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return strings.Join(messages, "\n"), nil
}
