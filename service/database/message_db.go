package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func (db *appdbimpl) SendMessage(senderId string, chatId string, text string, replyTo *string) (*Message, error) {
	now := time.Now()

	message := Message{
		MessageId:   generateID(),
		ChatId:      chatId,
		SenderId:    senderId,
		Text:        text,
		State:       StateSent,
		ContentType: TypeText,
		ReplyTo:     replyTo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tx, err := db.c.Begin()
	if err != nil {
		return nil, fmt.Errorf("Errore avvio transazione: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	result, err := tx.Exec(`INSERT INTO messages 
		(message_id, chat_id, sender_id, text, state, content_type, photo_url, reply_to, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		message.MessageId, message.ChatId, message.SenderId, message.Text, message.State, message.ContentType, message.PhotoUrl, message.ReplyTo, now, now)
	if err != nil {
		return nil, fmt.Errorf("Errore durante la creazione del messaggio: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("Errore nel recupero delle righe modificate: %w", err)
	}
	if rows == 0 {
		return nil, errors.New("Riga del messaggio non creata")
	}

	_, err = tx.Exec("UPDATE chats SET modified_at = ? WHERE chat_id = ?", now, message.ChatId)
	if err != nil {
		return nil, fmt.Errorf("Errore durante la modifica del timestamp della chat: %w", err)
	}

	// Se la chat è stata cancellata dall'utente, quell'utente rivedrà di nuovo la chat
	_, _ = tx.Exec("DELETE FROM deleted_chats WHERE chat_id = ?", message.ChatId)

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("Errore commit transazione: %w", err)
	}

	return &message, nil
}

func (db *appdbimpl) SendImageMessage(senderId string, chatId string, text string, photoUrl string, replyTo *string) (*Message, error) {
	now := time.Now()

	message := Message{
		MessageId:   generateID(),
		ChatId:      chatId,
		SenderId:    senderId,
		Text:        text,
		State:       StateSent,
		ContentType: TypeImage,
		PhotoUrl:    photoUrl,
		ReplyTo:     replyTo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tx, err := db.c.Begin()
	if err != nil {
		return nil, fmt.Errorf("Errore avvio transazione: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	result, err := tx.Exec(`INSERT INTO messages 
		(message_id, chat_id, sender_id, text, state, content_type, photo_url, reply_to, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		message.MessageId, message.ChatId, message.SenderId, message.Text, message.State, message.ContentType, message.PhotoUrl, message.ReplyTo, now, now)
	if err != nil {
		return nil, fmt.Errorf("Errore durante la creazione del messaggio: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("Errore nel recupero delle righe modificate: %w", err)
	}
	if rows == 0 {
		return nil, errors.New("Riga del messaggio non creata")
	}

	_, err = tx.Exec("UPDATE chats SET modified_at = ? WHERE chat_id = ?", now, message.ChatId)
	if err != nil {
		return nil, fmt.Errorf("Errore durante la modifica del timestamp della chat: %w", err)
	}

	// Se la chat è stata cancellata dall'utente, quell'utente rivedrà di nuovo la chat
	_, _ = tx.Exec("DELETE FROM deleted_chats WHERE chat_id = ?", message.ChatId)

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("Errore commit transazione: %w", err)
	}

	return &message, nil
}

func (db *appdbimpl) UpdateMessage(messageId string, text string, senderId string) (*Message, error) {
	if len(text) < 1 || len(text) > 1000 {
		return nil, errors.New("Il testo deve avere lunghezza tra 1 e 1000 caratteri")
	}
	now := time.Now()
	result, err := db.c.Exec("UPDATE messages SET text = ?, updated_at = ? WHERE message_id = ? AND sender_id = ?", text, now, messageId, senderId)
	if err != nil {
		return nil, fmt.Errorf("Errore durante la modifica del messaggio: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, errors.New("Messaggio non trovato o non autorizzato")
	}

	var msg Message
	var dbContentType string
	var dbState string
	err = db.c.QueryRow("SELECT message_id, chat_id, sender_id, text, state, content_type, COALESCE(photo_url, ''), created_at, updated_at FROM messages WHERE message_id = ?", messageId).Scan(&msg.MessageId, &msg.ChatId, &msg.SenderId, &msg.Text, &dbState, &dbContentType, &msg.PhotoUrl, &msg.CreatedAt, &msg.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("Errore nel recupero del messaggio aggiornato: %w", err)
	}
	msg.State = MessageState(dbState)
	msg.ContentType = ContentType(dbContentType)

	return &msg, nil
}

func (db *appdbimpl) DeleteMessage(messageId string, userId string) error {

	result, err := db.c.Exec("DELETE FROM messages WHERE message_id = ? AND sender_id = ?", messageId, userId)
	if err != nil {
		return fmt.Errorf("Errore durante la cancellazione del messaggio: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Errore durante la cancellazione del messaggio; %w", err)
	}

	// 404
	if rows == 0 {
		return errors.New("Messaggio non trovato")
	}

	return nil
}

func (db *appdbimpl) UpdateMessageStatus(messageId string, newStatus MessageState) (time.Time, error) {
	var currentStatus string
	if newStatus != StateReceived && newStatus != StateRead {
		return time.Time{}, fmt.Errorf("stato non valido: %s", newStatus)
	}
	err := db.c.QueryRow("SELECT state FROM messages WHERE message_id = ?", messageId).Scan(&currentStatus)
	if errors.Is(err, sql.ErrNoRows) {
		// 404
		return time.Time{}, errors.New("messaggio non trovato")
	}
	if err != nil {
		return time.Time{}, fmt.Errorf("errore lettura stato: %w", err)
	}
	curState := MessageState(currentStatus)
	isValid := false
	if curState == newStatus {
		return time.Now(), nil
	}
	switch curState {
	case StateSent:
		if newStatus == StateReceived || newStatus == StateRead {
			isValid = true
		}
	case StateReceived:
		if newStatus == StateRead {
			isValid = true
		}
	case StateRead:
		isValid = false
	}

	if !isValid {
		return time.Time{}, fmt.Errorf("transizione di stato non valida: da %s a %s", curState, newStatus) // 400
	}

	now := time.Now()
	_, err = db.c.Exec("UPDATE messages SET state = ?, updated_at = ? WHERE message_id = ?", newStatus, now, messageId)
	if err != nil {
		return time.Time{}, fmt.Errorf("errore aggiornamento stato: %w", err)
	}
	return now, nil
}

func (db *appdbimpl) ForwardMessage(messageId string, targetChatId string, senderId string) (*Message, error) {
	checkTargetUser := false
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM chat_members WHERE chat_id = ? AND user_id = ?)", targetChatId, senderId).Scan(&checkTargetUser)
	if err != nil {
		return nil, err
	}
	if !checkTargetUser {
		return nil, errors.New("L'utente che invia il messaggio non fa parte della chat di destinazione")
	}
	var contentType string
	var text string
	var photoUrl string

	err = db.c.QueryRow("SELECT text, content_type, COALESCE(photo_url, '') FROM messages WHERE message_id = ?", messageId).Scan(&text, &contentType, &photoUrl)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("messaggio originale non trovato")
	}
	if err != nil {
		return nil, err
	}

	// Inizializzo il messaggio
	var newMessage Message
	newMessage.MessageId = generateID()
	newMessage.ChatId = targetChatId
	newMessage.SenderId = senderId
	newMessage.State = StateSent
	newMessage.ContentType = ContentType(contentType)
	newMessage.CreatedAt = time.Now()
	newMessage.UpdatedAt = time.Now()
	newMessage.Text = text
	newMessage.PhotoUrl = photoUrl

	_, err = db.c.Exec(`INSERT INTO messages 
		(message_id, chat_id, sender_id, text, state, content_type, photo_url, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		newMessage.MessageId, newMessage.ChatId, newMessage.SenderId, newMessage.Text, string(newMessage.State), string(newMessage.ContentType), newMessage.PhotoUrl, newMessage.CreatedAt, newMessage.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("Errore durante l'inoltro del messaggio: %w", err)
	}

	return &newMessage, nil
}
