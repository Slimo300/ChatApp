package mock

import (
	"time"

	"github.com/Slimo300/ChatApp/backend/src/communication"
	"github.com/Slimo300/ChatApp/backend/src/database"
	"github.com/Slimo300/ChatApp/backend/src/models"
)

func (m *MockDB) GetGroupMessages(groupID uint, offset, num int) (messages []models.Message, err error) {
	for _, message := range m.Messages {
		for _, member := range m.Members {
			if message.MemberID == member.ID && member.GroupID == groupID {
				message.Member = member
				messages = append(messages, message)
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
