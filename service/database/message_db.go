package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func (db *appdbimpl) SendMessage(message Message) (*Message,error) {
	if message.State == "" {
        message.State = StateSent
    }
    if message.ContentType == "" {
        message.ContentType = TypeText
    }
	message.MessageId = generateID()
	now := time.Now()
	message.CreatedAt = now
	message.UpdatedAt = now
	
	result, err := db.c.Exec(`INSERT INTO messages 
		(message_id, chat_id, sender_id, text, state, content_type ,photo_url, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`, 
		message.MessageId, message.ChatId, message.SenderId, message.Text, message.State, message.ContentType, message.PhotoUrl, now, now)
	if err != nil {
		return nil, fmt.Errorf("Errore durante la creazione del messaggio: %w", err)
	}

	rows, err := result.RowsAffected()
	if rows == 0 {
		return nil, errors.New("Riga del messaggio non creata")
	}
	
	result, err = db.c.Exec("UPDATE chats SET modified_at = ? WHERE chat_id = ?", now, message.ChatId)
	if err != nil {
		return nil, fmt.Errorf("Errore durante la modifica del timestamp della chat: %w", err)
	}
	return &message, nil 
}

func (db *appdbimpl) UpdateMessage(message Message, text string) (*Message, error) {
	if len(text) < 1 || len(text) > 1000 {
		return nil, errors.New("Il testo deve avere lunghezza tra 1 e 1000 caratteri")
	}
	now := time.Now()
	_, err := db.c.Exec("UPDATE messages SET text = ?, updated_at = ? WHERE messageId = ?",text, now, message.MessageId)
	if err != nil {
		return nil, fmt.Errorf("Errore durante la modifica del messaggio: %w", err)
	}
	message.UpdatedAt = now
	return &message, nil
}

func (db *appdbimpl) DeleteMessage(messageId string, userId string) (string,error) {
	var userExists bool
	
	err := db.c.QueryRow("SELECT 1 FROM messages WHERE sender_id = ? AND message_id = ?", userId, messageId).Scan(&userExists)

	if errors.Is(err, sql.ErrNoRows) {
        return "", errors.New("Impossibile cancellare il messaggio, utente non autorizzato o messaggio non trovato") 
    }
	if err != nil {
		return "", fmt.Errorf("Errore durante la verifica del messaggio: %w", err)
	}
	
	result, err := db.c.Exec("DELETE FROM messages WHERE message_id = ? AND sender_id = ?", messageId, userId)
	if err != nil {
		return "", fmt.Errorf("Errore durante la cancellazione del messaggio: %w", err)
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("Errore durante la cancellazione del messaggio; %w", err)
	}
	
	if rows == 0 {
		// 404
		return "", errors.New("Messaggio non trovato")
	}
	
	return messageId, nil
}

func (db *appdbimpl) UpdateMessageStatus(messageId string, newStatus MessageState) (time.Time, error) {
    var currentStatus string
	if newStatus != StateReceived && newStatus != StateRead {
        return time.Time{}, fmt.Errorf("stato non valido: %s", string(newStatus))
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
        return time.Time{}, fmt.Errorf("transizione di stato non valida: da %s a %s", string(curState), string(newStatus)) // 400
    }

    now := time.Now()
    // Aggiornamento updated_at 
    _, err = db.c.Exec("UPDATE messages SET state = ?, updated_at = ? WHERE message_id = ?", string(newStatus), now, string(messageId))
    if err != nil {
        return time.Time{}, fmt.Errorf("errore aggiornamento stato: %w", err)
    }
    return now, nil
}

func (db *appdbimpl) ForwardMessage(messageId string, targetChatId string, senderId string) (*Message,error) {
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