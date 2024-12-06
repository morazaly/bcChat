package repository

import (
	"database/sql"
	"fmt"
	"time"

	"bccChat/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) SaveMessage(message models.SendMessage, username string) error {
	query := `INSERT INTO messages (sender_id, room_id, message, timestamp) VALUES (?, ?, ?, sysdate())`
	_, err := r.db.Exec(query, username, message.RoomID, message.Content)

	return err
}

func (r *MessageRepository) GetHistory(roomID string) ([]models.Message, error) {
	query := `SELECT id, sender_id, room_id, message, timestamp FROM messages WHERE room_id = ? ORDER BY timestamp ASC`

	rows, err := r.db.Query(query, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		var id int
		var senderID string
		var timestamp time.Time
		if err := rows.Scan(&id, &senderID, &msg.RoomID, &msg.Content, &timestamp); err != nil {
			return nil, err
		}
		fmt.Println("senderID", senderID)
		msg.Username = senderID
		msg.Timestamp = timestamp.Format(time.RFC3339)
		messages = append(messages, msg)
	}
	return messages, nil
}
