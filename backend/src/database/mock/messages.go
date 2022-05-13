package mock

import (
	"errors"
	"time"

	"github.com/Slimo300/ChatApp/backend/src/communication"
	"github.com/Slimo300/ChatApp/backend/src/database"
	"github.com/Slimo300/ChatApp/backend/src/models"
)

func (m *MockDB) GetGroupMessages(id_user, id_group uint, offset, num int) ([]communication.Message, error) {
	var messages []communication.Message

	for _, member := range m.Members {
		if member.GroupID == id_group && member.UserID == id_user {
			break
		}
		return nil, errors.New("User cannot request from this group")
	}

	for _, message := range m.Messages {
		for _, member := range m.Members {
			if message.MemberID == member.ID && member.GroupID == id_group {
				messages = append(messages, communication.Message{
					Group:   uint64(id_group),
					Member:  uint64(member.ID),
					Message: message.Text,
					Nick:    member.Nick,
					When:    message.Posted.Format(database.TIME_FORMAT),
				})
			}
		}
	}

	return messages, nil
}

func (m *MockDB) AddMessage(msg communication.Message) (communication.Message, error) {
	msgTime := time.Now()
	m.Messages = append(m.Messages, models.Message{
		ID:       uint(len(m.Messages) + 1),
		Posted:   msgTime,
		Text:     msg.Message,
		MemberID: uint(msg.Member),
	})

	msg.When = msgTime.Format(database.TIME_FORMAT)
	return msg, nil
}
