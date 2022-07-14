package handlers

import (
	"net/http"
	"time"

	"github.com/Slimo300/ChatApp/backend/lib/communication"
	"github.com/Slimo300/ChatApp/backend/lib/database"
	"github.com/Slimo300/ChatApp/backend/lib/storage"
	"github.com/Slimo300/ChatApp/backend/lib/ws"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type Server struct {
	DB           database.DBlayer
	Storage      storage.StorageLayer
	Hub          ws.HubInterface
	actionChan   chan<- *communication.Action
	messageChan  <-chan *communication.Message
	secret       string
	domain       string
	maxBodyBytes int
}

func NewServer(db database.DBlayer, storage storage.StorageLayer) *Server {
	actionChan := make(chan *communication.Action)
	messageChan := make(chan *communication.Message)
	return &Server{
		DB:           db,
		Storage:      storage,
		secret:       "wołowina",
		domain:       "localhost",
		actionChan:   actionChan,
		messageChan:  messageChan,
		maxBodyBytes: 4194304,
		Hub:          ws.NewHub(messageChan, actionChan),
	}
}

func NewServerWithMockHub(db database.DBlayer, storage storage.StorageLayer) *Server {
	actionChan := make(chan *communication.Action)
	messageChan := make(chan *communication.Message)
	return &Server{
		DB:           db,
		Storage:      storage,
		secret:       "wołowina",
		domain:       "localhost",
		actionChan:   actionChan,
		messageChan:  messageChan,
		maxBodyBytes: 4194304,
		Hub:          ws.NewMockHub(actionChan),
	}
}

func (s *Server) RunHub() {
	go s.ListenToHub()
	s.Hub.Run()
}

func (s *Server) ListenToHub() {
	var msg *communication.Message
	for {
		select {
		case msg = <-s.messageChan:
			when, err := time.Parse(communication.TIME_FORMAT, msg.When)
			if err != nil {
				panic(err.Error())
			}
			if err := s.DB.AddMessage(msg.Member, msg.Message, when); err != nil {
				panic("Panicked while adding message")
			}
		}
	}
}

func (s *Server) CreateSignedToken(iss string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    iss,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *Server) MustAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("jwt")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
			return
		}
		token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{},
			func(t *jwt.Token) (interface{}, error) {
				return []byte(s.secret), nil
			})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
			return
		}
		userID := token.Claims.(*jwt.StandardClaims).Issuer
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
			return
		}
		c.Set("userID", userID)
		c.Next()
	}
}

// middleware for checking database connection
func (s *Server) CheckDatabase() gin.HandlerFunc {
	return func(c *gin.Context) {
		if s.DB == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": ErrNoDatabase.Error()})
			return
		}
		c.Next()
	}
}
