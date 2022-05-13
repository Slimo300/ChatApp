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
	DB       database.DBlayer
	Hub      *ws.Hub
	CommChan chan<- *communication.Action
	secret   string
	domain   string
}

func NewServer(db database.DBlayer, ch chan *communication.Action) *Server {
	return &Server{
		DB:       db,
		secret:   "wołowina",
		domain:   "localhost",
		CommChan: ch,
		Hub:      ws.NewHub(db, ch),
	}
}

func (s *Server) RunHub() {
	s.Hub.Run()
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
