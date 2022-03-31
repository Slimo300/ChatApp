package database

import (
	"time"

	"github.com/Slimo300/ChatApp/backend/src/communication"
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
		When:    message.Posted.Format(TIME_FORMAT),
	}, nil

}

// GetGroupMessages gets last messages sent to group, offset is a variable stating how many of messages it should omit (in case some are already loaded)
func (db *Database) GetGroupMessages(id_user, id_group, offset uint) ([]communication.Message, error) {

	if err := db.Where(models.Member{GroupID: id_group, UserID: id_user}).First(&models.Member{}).Error; err != nil {
		return nil, err
	}

	var messages []models.Message
	selection := db.Joins("Member", db.Where(&models.Member{GroupID: id_group})).Order("posted desc").Offset(int(offset)*15).Limit(MESSAGE_LIMIT).Find(&messages, "`Member`.`group_id` = ?", id_group)
	if selection.Error != nil {
		return nil, selection.Error
	}
	var sendMessages []communication.Message

	for i := len(messages) - 1; i >= 0; i-- {
		sendMessages = append(sendMessages, communication.Message{
			Group:   uint64(messages[i].Member.GroupID),
			Member:  uint64(messages[i].MemberID),
			Nick:    messages[i].Member.Nick,
			When:    messages[i].Posted.Format(TIME_FORMAT),
			Message: messages[i].Text,
		})
	}

	return sendMessages, nil
}
