package database

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

func (db *appdbimpl) CreateGroup(groupName string, groupDescription string, photoUrl string, creatorId string) (*Group, error) {
	if len(groupName) < 1 || len(groupName) > 30 {
		return nil, errors.New("Il nome del gruppo deve avere tra 1 e 30 caratteri")
	}
	if len(groupDescription) > 50 {
		return nil, errors.New("La descrizione non può superare i 50 caratteri")
	}
	// Controllo che l'utente esista
	checkUser := false
	err := db.c.QueryRow("SELECT 1 FROM users WHERE identifier = ?", creatorId).Scan(&checkUser)
	if err != nil {
		return nil, err
	}
	if !checkUser {
		return nil, errors.New("utente creatore non trovato")
	}

	// Creo il gruppo
	groupId := generateID()
	now := time.Now()
	_, err = db.c.Exec(`INSERT INTO groups
		(group_id, group_name, group_description, photo_url, created_at) VALUES (?, ?, ?, ?, ?)`,
		groupId, groupName, groupDescription, photoUrl, now)
	if err != nil {
		return nil, fmt.Errorf("Errore durante l'inserimento del gruppo: %w", err)
	}

	// Inserisco tra i membri l'utente che ha creato il gruppo
	err = db.AddToGroup(groupId, creatorId)
	if err != nil {
		return nil, fmt.Errorf("Errore durante l'aggiunta dell'utente nel gruppo: %w", err)
	}

	group := Group{
		GroupId:          groupId,
		GroupName:        groupName,
		GroupDescription: groupDescription,
		PhotoUrl:         photoUrl,
		CreatedAt:        now,
		Members:          []string{creatorId},
	}
	return &group, nil
}

func (db *appdbimpl) GetGroup(groupId string) (*Group, error) {
	// Controllo che il gruppo esista
	var groupExists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM groups WHERE group_id = ?)", groupId).Scan(&groupExists)
	if err != nil {
		return nil, err
	}
	if !groupExists {
		return nil, errors.New("gruppo non trovato")
	}
	// Query Gruppo
	var g Group
	err = db.c.QueryRow(`SELECT group_id, group_name, group_description, photo_url, created_at 
		FROM groups WHERE group_id = ?`, groupId).Scan(&g.GroupId, &g.GroupName, &g.GroupDescription, &g.PhotoUrl, &g.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("errore lettura gruppo: %w", err)
	}

	// Query Membri
	m := []string{}
	result, err := db.c.Query("SELECT user_id FROM chat_members WHERE group_id = ?", groupId)
	if err != nil {
		return nil, fmt.Errorf("errore lettura membri: %w", err)
	}
	defer result.Close()

	for result.Next() {
		var userId string
		err := result.Scan(&userId)
		if err != nil {
			return nil, fmt.Errorf("Errore lettura ID di un membro del gruppo: %w", err)
		}
		m = append(m, userId)
	}
	g.Members = m
	return &g, nil
}

func (db *appdbimpl) SetGroupName(groupId string, userId string, newName string) error {
	// Controllo che il gruppo esista
	var groupExists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM groups WHERE group_id = ?)", groupId).Scan(&groupExists)
	if err != nil {
		return err
	}
	if !groupExists {
		return errors.New("gruppo non trovato")
	}

	// Controllo se l'utente è admin del gruppo
	var isAdmin bool
	err = db.c.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM chat_members cm
			JOIN groups g ON cm.chat_id = g.chat_id
			WHERE g.group_id = ? AND cm.user_id = ? AND cm.role = 'admin'
		)`, groupId, userId).Scan(&isAdmin)
	if err != nil {
		return fmt.Errorf("Errore durante la verifica dei permessi: %w", err)
	}
	if !isAdmin {
		return errors.New("Utente non autorizzato a modificare il nome del gruppo")
	}

	if len(newName) < 1 || len(newName) > 30 {
		return errors.New("La lunghezza del nome del gruppo deve essere tra 1 e 30 caratteri")
	}

	result, err := db.c.Exec("UPDATE groups SET group_name = ? WHERE group_id = ?", newName, groupId)
	if err != nil {
		return fmt.Errorf("Errore nell'aggiornamento del nome del gruppo: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		// 404
		return errors.New("Gruppo non trovato")
	}

	return nil
}

func (db *appdbimpl) SetGroupPhoto(groupId string, userId string, photoUrl string) error {
	if len(photoUrl) < 5 || len(photoUrl) > 50 {
		return errors.New("La lunghezza dell'URL della foto deve essere tra i 5 e i 50 caratteri")
	}

	// Controllo se l'utente è admin del gruppo
	var isAdmin bool
	err := db.c.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM chat_members cm
			JOIN groups g ON cm.chat_id = g.chat_id
			WHERE g.group_id = ? AND cm.user_id = ? AND cm.role = 'admin'
		)
	`, groupId, userId).Scan(&isAdmin)
	if err != nil {
		return fmt.Errorf("Errore durante la verifica dei permessi: %w", err)
	}
	if !isAdmin {
		return errors.New("Utente non autorizzato a modificare la foto del gruppo")
	}

	_, err = db.c.Exec("UPDATE groups SET photo_url = ? WHERE group_id = ?", photoUrl, groupId)
	if err != nil {
		return fmt.Errorf("Errore nell'aggiornamento della foto del gruppo: %w", err)
	}
	return nil
}

func (db *appdbimpl) AddToGroup(groupId string, userId string) error {
	// Controllo l'esistenza del gruppo e dell'utente
	var groupExists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM groups WHERE group_id = ?)", groupId).Scan(&groupExists)
	if err != nil {
		return err
	}
	if !groupExists {
		return errors.New("gruppo non trovato")
	}
	var userExists bool
	err = db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE identifier = ?)", userId).Scan(&userExists)
	if err != nil {
		return err
	}
	if !userExists {
		return errors.New("utente non trovato")
	}

	_, err = db.c.Exec("INSERT INTO group_members (group_id, user_id) VALUES (?, ?)", groupId, userId)
	if err != nil {
		// 409
		// Se l'errore contiene "UNIQUE constraint", significa che la coppia esiste già
		if strings.Contains(err.Error(), "UNIQUE") {
			return errors.New("l'utente è già membro del gruppo")
		}
		return fmt.Errorf("errore aggiunta membro: %w", err)
	}
	return nil
}

func (db *appdbimpl) LeaveGroup(groupId string, userId string) error {

	var groupExists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM groups WHERE group_id = ?)", groupId).Scan(&groupExists)
	if err != nil {
		return err
	}
	if !groupExists {
		return errors.New("gruppo non trovato")
	}
	var userExists bool
	err = db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE identifier = ?)", userId).Scan(&userExists)
	if err != nil {
		return err
	}
	if !userExists {
		return errors.New("utente non trovato")
	}

	_, err = db.c.Exec("DELETE FROM group_members WHERE group_id = ? AND member_id = ?", groupId, userId)
	if err != nil {
		return errors.New("Utente non trovato oppure gruppo inesistente")
	}
	return nil
}
