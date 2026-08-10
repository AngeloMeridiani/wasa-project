package database

import (
	"database/sql"
	"fmt"
	"time"
)

func (db *appdbimpl) CreateChat(userId1 string, userId2 string) (*Chat, error) {

	chatId := generateID()
	createdAt := time.Now()
	modifiedAt := time.Now()

	// Check if a DM chat already exists
	var existingChatId string
	var extCreatedAt, extModifiedAt time.Time
	var checkErr error

	if userId1 == userId2 {
		checkErr = db.c.QueryRow(`
			SELECT c.chat_id, c.created_at, c.modified_at
			FROM chats c
			WHERE c.type = 'dm'
			AND (SELECT COUNT(*) FROM chat_members WHERE chat_id = c.chat_id) = 1
			AND (SELECT user_id FROM chat_members WHERE chat_id = c.chat_id LIMIT 1) = ?
		`, userId1).Scan(&existingChatId, &extCreatedAt, &extModifiedAt)
	} else {
		checkErr = db.c.QueryRow(`
			SELECT c.chat_id, c.created_at, c.modified_at
			FROM chats c
			JOIN chat_members cm1 ON c.chat_id = cm1.chat_id AND cm1.user_id = ?
			JOIN chat_members cm2 ON c.chat_id = cm2.chat_id AND cm2.user_id = ?
			WHERE c.type = 'dm'
		`, userId1, userId2).Scan(&existingChatId, &extCreatedAt, &extModifiedAt)
	}

	if checkErr == nil && existingChatId != "" {
		var title string
		if userId1 == userId2 {
			title = "Messaggi Salvati"
		} else {
			err := db.c.QueryRow("SELECT name FROM users WHERE identifier = ?", userId2).Scan(&title)
			if err != nil {
				title = existingChatId
			}
		}
		return &Chat{
			ChatId:     existingChatId,
			Title:      title,
			CreatedAt:  extCreatedAt,
			ModifiedAt: extModifiedAt,
		}, nil
	}

	// Iniziamo una transazione
	tx, err := db.c.Begin()
	if err != nil {
		return nil, fmt.Errorf("Errore avvio transazione: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// Creiamo la chat nella tabella chats
	_, err = tx.Exec("INSERT INTO chats (chat_id, created_at, modified_at, type) VALUES (?, ?, ?, ?)", chatId, createdAt, modifiedAt, "dm")
	if err != nil {
		return nil, fmt.Errorf("errore creazione chat: %w", err)
	}

	if userId1 == userId2 {
		_, err = tx.Exec("INSERT INTO chat_members (chat_id, user_id) VALUES (?, ?)", chatId, userId1)
	} else {
		_, err = tx.Exec("INSERT INTO chat_members (chat_id, user_id) VALUES (?, ?), (?, ?)", chatId, userId1, chatId, userId2)
	}

	if err != nil {
		return nil, fmt.Errorf("errore inserimento membri: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("errore commit: %w", err)
	}

	var title string
	if userId1 == userId2 {
		title = "Utente senza nome"
	} else {
		err = db.c.QueryRow("SELECT name FROM users WHERE identifier = ?", userId2).Scan(&title)
		if err != nil {
			title = chatId // fallback
		}
	}

	newChat := &Chat{
		ChatId:     chatId,
		Title:      title,
		CreatedAt:  createdAt,
		ModifiedAt: modifiedAt,
	}

	return newChat, nil
}

func (db *appdbimpl) GetCurrentChat(chatId string, userId string, limit int) (*Chat, []Message, error) {
	chatQuery := `
		SELECT 
			c.chat_id, 
			c.created_at, 
			c.modified_at,
			CASE 
				WHEN c.type = 'group' THEN g.group_name
				WHEN u.name IS NOT NULL THEN u.name
				ELSE 'Messaggi Salvati'
			END as title
		FROM chats c
		JOIN chat_members cm1 ON c.chat_id = cm1.chat_id AND cm1.user_id = ?
		LEFT JOIN groups g ON c.chat_id = g.chat_id AND c.type = 'group'
		LEFT JOIN chat_members cm2 ON c.chat_id = cm2.chat_id AND cm2.user_id != ? AND c.type = 'dm'
		LEFT JOIN users u ON cm2.user_id = u.identifier
		WHERE c.chat_id = ?
	`

	var chat Chat
	err := db.c.QueryRow(chatQuery, userId, userId, chatId).Scan(
		&chat.ChatId, &chat.CreatedAt, &chat.ModifiedAt, &chat.Title,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("chat non trovata o non autorizzato: %w", err)
	}

	msgQuery := `
		SELECT * FROM (
			SELECT message_id, sender_id, text, state, content_type, COALESCE(photo_url, ''), created_at, updated_at
			FROM messages
			WHERE chat_id = ?
			ORDER BY created_at DESC
			LIMIT ?
		) AS sub
		ORDER BY created_at ASC
	`
	rows, err := db.c.Query(msgQuery, chatId, limit)
	if err != nil {
		return nil, nil, fmt.Errorf("errore lettura messaggi: %w", err)
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var m Message
		var dbState, dbContentType string
		err := rows.Scan(
			&m.MessageId, &m.SenderId, &m.Text, &dbState, &dbContentType, &m.PhotoUrl, &m.CreatedAt, &m.UpdatedAt,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("errore scan messaggio: %w", err)
		}
		m.ChatId = chatId
		m.State = MessageState(dbState)
		m.ContentType = ContentType(dbContentType)
		messages = append(messages, m)
	}

	if err = rows.Err(); err != nil {
		return nil, nil, err
	}

	return &chat, messages, nil
}

func (db *appdbimpl) GetMyChats(userId string) ([]Chat, error) {
	query := `
		SELECT 
			c.chat_id, 
			c.created_at, 
			c.modified_at,
			CASE 
				WHEN c.type = 'group' THEN g.group_name
				WHEN u.name IS NOT NULL THEN u.name
				ELSE 'Messaggi Salvati'
			END as title,
			m.message_id,
			m.sender_id,
			m.text,
			m.state,
			m.content_type,
			m.photo_url,
			m.created_at as msg_created_at,
			m.updated_at as msg_updated_at
		FROM chats c
		JOIN chat_members cm1 ON c.chat_id = cm1.chat_id AND cm1.user_id = ?
		LEFT JOIN groups g ON c.chat_id = g.chat_id AND c.type = 'group'
		LEFT JOIN chat_members cm2 ON c.chat_id = cm2.chat_id AND cm2.user_id != ? AND c.type = 'dm'
		LEFT JOIN users u ON cm2.user_id = u.identifier
		LEFT JOIN (
			SELECT m1.*
			FROM messages m1
			JOIN (
				SELECT chat_id, MAX(created_at) as max_created_at
				FROM messages
				GROUP BY chat_id
			) m2 ON m1.chat_id = m2.chat_id AND m1.created_at = m2.max_created_at
			GROUP BY m1.chat_id 
		) m ON c.chat_id = m.chat_id
		WHERE c.chat_id NOT IN (
			SELECT chat_id FROM deleted_chats WHERE user_id = ?
		)
		ORDER BY c.modified_at DESC
	`

	rows, err := db.c.Query(query, userId, userId, userId)
	if err != nil {
		return nil, fmt.Errorf("errore nella query GetMyChats: %w", err)
	}
	defer rows.Close()

	var chats []Chat
	for rows.Next() {
		var chat Chat
		var title sql.NullString

		var msgId, senderId, text, state, contentType, photoUrl sql.NullString
		var msgCreatedAt, msgUpdatedAt sql.NullTime

		err := rows.Scan(
			&chat.ChatId, &chat.CreatedAt, &chat.ModifiedAt, &title,
			&msgId, &senderId, &text, &state, &contentType, &photoUrl, &msgCreatedAt, &msgUpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("errore nello scan della chat: %w", err)
		}
		if title.Valid {
			chat.Title = title.String
		}

		if msgId.Valid {
			chat.LastMessage = &Message{
				MessageId:   msgId.String,
				ChatId:      chat.ChatId,
				SenderId:    senderId.String,
				Text:        text.String,
				State:       MessageState(state.String),
				ContentType: ContentType(contentType.String),
			}
			if photoUrl.Valid {
				chat.LastMessage.PhotoUrl = photoUrl.String
			}
			if msgCreatedAt.Valid {
				chat.LastMessage.CreatedAt = msgCreatedAt.Time
			}
			if msgUpdatedAt.Valid {
				chat.LastMessage.UpdatedAt = msgUpdatedAt.Time
			}
		}

		chats = append(chats, chat)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if chats == nil {
		chats = make([]Chat, 0)
	}

	return chats, nil
}

func (db *appdbimpl) DeleteChat(chatId string, userId string) error {
	var chatType string
	err := db.c.QueryRow("SELECT type FROM chats WHERE chat_id = ?", chatId).Scan(&chatType)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("chat non trovata")
		}
		return fmt.Errorf("errore ricerca chat: %w", err)
	}

	if chatType == "group" {
		return db.LeaveGroup(chatId, userId)
	}

	// For DMs, soft delete by inserting into deleted_chats
	_, err = db.c.Exec("INSERT OR IGNORE INTO deleted_chats (chat_id, user_id) VALUES (?, ?)", chatId, userId)
	if err != nil {
		return fmt.Errorf("errore eliminazione chat: %w", err)
	}
	return nil
}
