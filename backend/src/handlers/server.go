package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Slimo300/ChatApp/backend/src/communication"
	"github.com/Slimo300/ChatApp/backend/src/database"
	"github.com/Slimo300/ChatApp/backend/src/ws"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type Server struct {
	DB          database.DBlayer
	Hub         *ws.Hub
	sendHubChan chan<- *communication.Action
	rcvHubChan  <-chan communication.Message
	secret      string
	domain      string
}

func NewServer(db database.DBlayer) *Server {
	sendHubChan := make(chan *communication.Action)
	rcvHubChan := make(chan communication.Message)
	return &Server{
		DB:          db,
		secret:      "wołowina",
		domain:      "localhost",
		sendHubChan: sendHubChan,
		rcvHubChan:  rcvHubChan,
		Hub:         ws.NewHub(rcvHubChan, sendHubChan),
	}
}

func (s *Server) ListenToHub() {
	var msg communication.Message
	select {
	case msg = <-s.rcvHubChan:
		if err := s.DB.AddMessage(uint(msg.Member), msg.Message); err != nil {
			panic("Panicced while adding message")
		}
	}
}

func (s *Server) RunHub() {
	go s.Hub.Run()
	s.ListenToHub()
}

func (s *Server) CreateSignedToken(iss int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(iss),
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
		id, err := strconv.Atoi(token.Claims.(*jwt.StandardClaims).Issuer)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
			return
		}
		c.Set("userID", id)
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
