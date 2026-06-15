package database

import (
	"errors"
	"fmt"
	"time"
)

func (db *appdbimpl) CreateComment(senderId string, messageId string, emoji string) (*Comment, error) {
	commentId := generateID()
	now := time.Now()
	comment := Comment{
		CommentId: commentId,
		MessageId: messageId,
		SenderId:  senderId,
		Emoji:     emoji,
		CreatedAt: now,
	}

	result, err := db.c.Exec(`INSERT INTO comments 
		(comment_id, message_id, sender_id, emoji, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`,
		comment.CommentId, comment.MessageId, comment.SenderId, comment.Emoji, now, now)

	if err != nil {
		return nil, fmt.Errorf("Errore durante la creazione del commento: %w", err)
	}
	rows, err := result.RowsAffected()
	if rows == 0 || err != nil {
		return nil, errors.New("Riga del commento non creata")
	}
	return &comment, nil
}

func (db *appdbimpl) DeleteComment(commentId string, userId string) error {

	result, err := db.c.Exec("DELETE from comments WHERE comment_id = ? AND sender_id = ?", commentId, userId)
	if err != nil {
		return fmt.Errorf("Errore durante la cancellazione del commento: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Errore durante la cancellazione del commento: %w", err)
	}

	if rows == 0 {
		return errors.New("Riga del commento non trovata")
	}
	return nil
}
