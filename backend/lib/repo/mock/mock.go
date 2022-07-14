package mock

import (
	"time"

	"github.com/Slimo300/ChatApp/backend/lib/repo"
)

type TokenInfo struct {
	Created    time.Time
	Expiration time.Duration
	Value      repo.TokenValue
}

type mockTokenRepository map[string]TokenInfo

func NewMockTokenRepository(host, port, password string) mockTokenRepository {
	return make(mockTokenRepository)
}

func (mock mockTokenRepository) SaveToken(token string, expiration time.Duration) error {
	mock[token] = TokenInfo{
		Created:    time.Now(),
		Expiration: expiration,
		Value:      repo.TOKEN_VALID,
	}
	return nil
}

// func (mock mockTokenRepository) IsTokenValid(userID, tokenID string) (bool, error) {
// 	pattern := fmt.Sprintf("%s:*%s", userID, tokenID)
// 	log.Println(pattern)

// 	keys, err := rdb.Keys(pattern).Result()
// 	if err != nil {
// 		return false, err
// 	}
// 	if len(keys) == 0 {
// 		return false, repo.TokenNotFoundError
// 	}
// 	if len(keys) > 1 {
// 		return false, repo.TooManyTokensFoundError
// 	}

// 	res, err := rdb.Get(keys[0]).Result()
// 	if err != nil {
// 		return false, err
// 	}
// 	if res != repo.TOKEN_VALID {
// 		if res == repo.TOKEN_BLACKLISTED {
// 			return false, repo.TokenBlacklistedError
// 		}
// 		return false, errors.New("Unexpected token value")
// 	}

// 	return true, nil
// }

// func (rdb *redisTokenRepository) InvalidateTokens(userID, tokenID string) error {
// 	t := tokenID
// 	for {
// 		key := fmt.Sprintf("%s:%s:*", userID, t)

// 		keys, err := rdb.Keys(key).Result()
// 		if err != nil {
// 			return err
// 		}
// 		if len(keys) == 0 {
// 			break
// 		}
// 		if len(keys) > 1 {
// 			return errors.New(fmt.Sprint("Too many tokens: ", keys))
// 		}
// 		if err := rdb.Do("set", keys[0], repo.TOKEN_BLACKLISTED, "keepttl").Err(); err != nil {
// 			return err
// 		}
// 		t = strings.Split(keys[0], ":")[2]
// 	}

// 	return nil
// }

// func (rdb *redisTokenRepository) InvalidateToken(userID, tokenID string) error {
// 	key := fmt.Sprintf("%s:*%s", userID, tokenID)

// 	keys, err := rdb.Keys(key).Result()
// 	if err != nil {
// 		return err
// 	}
// 	if len(keys) != 1 {
// 		return errors.New(fmt.Sprint("Wrong number of tokens: ", keys))
// 	}
// 	if err := rdb.Do("set", keys[0], repo.TOKEN_BLACKLISTED, "keepttl").Err(); err != nil {
// 		return err
// 	}
// 	return nil
// }
