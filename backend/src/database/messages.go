package database

import "github.com/Slimo300/ChatApp/backend/src/models"

func (db *Database) AddMessage(msg Message) error {
	message := models.Message{
		Posted:   msg.When,
		Text:     msg.Message,
		MemberID: uint(msg.Member),
	}

	if err := db.Create(&message).Error; err != nil {
		return err
	}

	return nil

}
