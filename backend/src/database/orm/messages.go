package orm

import (
	"time"

	"github.com/Slimo300/ChatApp/backend/src/communication"
	"github.com/Slimo300/ChatApp/backend/src/database"
	"github.com/Slimo300/ChatApp/backend/src/models"
)

func (db *Database) AddMessage(msg communication.Message) (communication.Message, error) {
	message := models.Message{
		Posted:   time.Now(),
		Text:     msg.Message,
		MemberID: uint(msg.Member),
	}

	if err := db.Create(&message).Error; err != nil {
		return communication.Message{}, err
	}

	var member models.Member
	if err := db.First(&member, message.MemberID).Error; err != nil {
		return communication.Message{}, err
	}

	return communication.Message{
		Group:   uint64(member.GroupID),
		Member:  uint64(message.MemberID),
		Message: message.Text,
		Nick:    message.Member.Nick,
		When:    message.Posted.Format(database.TIME_FORMAT),
	}, nil

}

func (db *Database) GetGroupMessages(groupID uint, offset, num int) (messages []models.Message, err error) {
	return messages, db.Joins("Member", db.Where(&models.Member{GroupID: groupID})).Order("posted desc").Offset(offset).Limit(num).
		Find(&messages, "`Member`.`group_id` = ?", groupID).Error
}
