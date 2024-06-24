package pq

import (
	"database/sql"
	"fmt"

	proto "chat/proto"
)

type ChatDB struct {
	db *sql.DB
}

func NewChatDB(db *sql.DB) *ChatDB {
	return &ChatDB{db: db}
}

func (c *ChatDB) SaveMessage(msg *proto.ChatMessage) error {
	query := `INSERT INTO messages (user, message, timestamp, ip_address) VALUES ($1, $2, $3, $4)`
	_, err := c.db.Exec(query, msg.User, msg.Message, msg.Timestamp, msg.IpAddress)
	if err != nil {
		return fmt.Errorf("failed to save message: %w", err)
	}
	return nil
}

func (c *ChatDB) LoadMessages() ([]*proto.ChatMessage, error) {
	query := `SELECT user, message, timestamp, ip_address FROM messages ORDER BY timestamp ASC`
	rows, err := c.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to load messages: %w", err)
	}
	defer rows.Close()

	var messages []*proto.ChatMessage
	for rows.Next() {
		var msg proto.ChatMessage
		if err := rows.Scan(&msg.User, &msg.Message, &msg.Timestamp, &msg.IpAddress); err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, &msg)
	}
	return messages, nil
}
