package mock

import (
	"time"

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

func (m *MockDB) AddMessage(memberID uint, text string, when time.Time) error {
	m.Messages = append(m.Messages, models.Message{
		ID:       uint(len(m.Messages) + 1),
		Posted:   when,
		Text:     text,
		MemberID: memberID,
	})

	return nil
}
