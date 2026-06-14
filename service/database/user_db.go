package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

func (db *appdbimpl) CreateUser(newName string) (*User, error) {

	trimmedName := strings.TrimSpace(newName)
	length := len(trimmedName)
	if length < 3 {
		return nil, errors.New("Il nome è troppo corto (minimo 3 caratteri)")
	}
	if length > 25 {
		return nil, errors.New("Il nome è troppo lungo (massimo 25 caratteri)")
	}

	userId := generateID()
	_, err := db.c.Exec("INSERT INTO users (identifier, name) VALUES (?, ?)", userId, newName)
	if err != nil {
		return nil, fmt.Errorf("errore creazione utente: %w", err)
	}

	user := &User{
		UserId: userId,
		Name:   newName,
	}

	return user, nil
}

func (db *appdbimpl) SetUserName(identifier string, newName string) error {

	trimmedName := strings.TrimSpace(newName)
	length := len(trimmedName)

	if length < 3 {
		return errors.New("Il nome è troppo corto (minimo 3 caratteri)")
	}
	if length > 25 {
		return errors.New("Il nome è troppo lungo (massimo 25 caratteri)")
	}

	result, err := db.c.Exec("UPDATE users SET name = ? WHERE identifier = ?", newName, identifier)
	if err != nil {
		// 400
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return fmt.Errorf("il nome '%s' è già in uso", trimmedName)
		}
		return fmt.Errorf("errore aggiornamento nome: %w", err)
	}

	righeToccate, _ := result.RowsAffected()
	// 404
	if righeToccate == 0 {
		return errors.New("utente non trovato")
	}

	return nil
}

func (db *appdbimpl) GetUser(userId string) (*User, error) {

	var user User
	// Uso COALESCE(photo_url, '') perché in Go le stringhe non possono essere NULL.
	// Se nel DB il campo è NULL, COALESCE lo trasforma in una stringa vuota "" per evitare errori.
	err := db.c.QueryRow("SELECT identifier, name, COALESCE(photo_url, '') FROM users WHERE identifier = ?", userId).Scan(&user.UserId, &user.Name, &user.UserPhoto)

	// 4. Gestione degli errori
	if errors.Is(err, sql.ErrNoRows) {
		// 404
		return nil, errors.New("utente non trovato")
	}
	if err != nil {
		return nil, fmt.Errorf("errore lettura utente: %w", err)
	}
	return &user, nil
}

func (db *appdbimpl) SetMyPhoto(id string, photoUrl string) error {
	length := len(photoUrl)
	if length < 5 || length > 50 {
		return errors.New("l'URL della foto deve essere tra 5 e 50 caratteri")
	}
	query := "UPDATE users SET photo_url = ? WHERE identifier = ?"

	result, err := db.c.Exec(query, photoUrl, id)
	if err != nil {
		return fmt.Errorf("errore aggiornamento foto: %w", err)
	}

	righeToccate, _ := result.RowsAffected()
	if righeToccate == 0 {
		// 404
		return errors.New("utente non trovato")
	}

	return nil
}
