package orm

import (
	"time"

	"github.com/Slimo300/ChatApp/backend/src/models"
)

func (db *Database) AddMessage(memberID uint, text string, when time.Time) error {
	message := models.Message{
		Text:     text,
		MemberID: memberID,
		Posted:   when,
	}

	if err := db.Create(&message).Error; err != nil {
		return err
	}

	return nil
}

func (db *Database) GetGroupMessages(groupID uint, offset, num int) (messages []models.Message, err error) {
	return messages, db.Joins("Member", db.Where(&models.Member{GroupID: groupID})).Order("posted desc").Offset(offset).Limit(num).
		Find(&messages, "`Member`.`group_id` = ?", groupID).Error
}
