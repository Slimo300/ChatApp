package mock

import (
	"errors"

	"github.com/Slimo300/ChatApp/backend/src/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (mock *MockDB) NewVerificationCode(userID uuid.UUID, code string) (verificationCode models.VerificationCode, err error) {
	verificationCode.UserID = userID
	verificationCode.ActivationCode = code
	mock.VerificationCodes = append(mock.VerificationCodes, verificationCode)
	return verificationCode, nil
}

func (mock *MockDB) VerifyCode(userID uuid.UUID, code string) error {

	for i, c := range mock.VerificationCodes {
		if c.UserID == userID {
			if code == c.ActivationCode {
				mock.VerificationCodes = append(mock.VerificationCodes[:i], mock.VerificationCodes[i+1:]...)
				for _, user := range mock.Users {
					if user.ID == userID {
						user.Activated = true
					}
				}
				return nil
			}
			return errors.New("Invalid code")
		}
	}
	return gorm.ErrRecordNotFound
}
