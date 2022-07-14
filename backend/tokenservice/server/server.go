package server

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Slimo300/ChatApp/backend/tokenservice/pb"
	"github.com/Slimo300/ChatApp/backend/tokenservice/repo"
	"github.com/golang-jwt/jwt"
	"golang.org/x/net/context"
)

type TokenService struct {
	*pb.UnimplementedTokenServiceServer
	repo                  repo.TokenRepository
	refreshTokenSecret    string
	accessTokenPrivateKey rsa.PrivateKey
	accessTokenDuration   time.Duration
	refreshTokenDuration  time.Duration
}

func NewTokenService(repo repo.TokenRepository, refreshSecret string, accessPrivKey rsa.PrivateKey,
	refreshDuration, accessDuration time.Duration) *TokenService {

	return &TokenService{
		repo:                  repo,
		refreshTokenSecret:    refreshSecret,
		accessTokenPrivateKey: accessPrivKey,
		refreshTokenDuration:  refreshDuration,
		accessTokenDuration:   accessDuration,
	}
}

func (srv *TokenService) NewPairFromUserID(ctx context.Context, userID *pb.UserID) (*pb.TokenPair, error) {
	log.Printf("Receive message body from client: %s", userID.ID)

	refreshData, err := srv.generateRefreshToken(userID.ID)
	if err != nil {
		log.Println("generate Refresh Error: ", err.Error())
		return &pb.TokenPair{
			Error: err.Error(),
		}, err
	}

	if err := srv.repo.SaveToken(fmt.Sprintf("%s:%s", userID.ID, refreshData.ID.String()), srv.refreshTokenDuration); err != nil {
		log.Println("saving refresh Error: ", err.Error())
		return &pb.TokenPair{
			Error: err.Error(),
		}, err
	}

	access, err := srv.generateAccessToken(userID.ID)
	if err != nil {
		log.Println("generate Access Error: ", err.Error())
		return &pb.TokenPair{
			Error: err.Error(),
		}, err
	}
	return &pb.TokenPair{
		AccessToken:  access,
		RefreshToken: refreshData.Token,
	}, nil
}

func (srv *TokenService) NewPairFromRefresh(ctx context.Context, refresh *pb.RefreshToken) (*pb.TokenPair, error) {
	log.Printf("Receive message body from client: %s", refresh.Token)

	token, err := jwt.ParseWithClaims(refresh.GetToken(), &jwt.StandardClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(srv.refreshTokenSecret), nil
		})
	if err != nil {
		log.Println("Parsing jwt error: ", err.Error())
		return &pb.TokenPair{
			Error: err.Error(),
		}, err
	}
	userID := token.Claims.(*jwt.StandardClaims).Subject
	tokenID := token.Claims.(*jwt.StandardClaims).Id
	log.Println(userID)
	log.Println(tokenID)

	ok, err := srv.repo.IsTokenValid(userID, tokenID)
	if err != nil {
		fmt.Println(err.Error())
		if errors.Is(err, repo.TokenBlacklistedError) {
			if err := srv.repo.InvalidateTokens(userID, tokenID); err != nil {
				panic(fmt.Sprint("Invalidating tokens failed: ", err.Error()))
			}
		}
		log.Println("Validating token error: ", err.Error())
		return &pb.TokenPair{
			Error: err.Error(),
		}, err
	}
	if !ok {
		log.Println("Token invalid")
		return &pb.TokenPair{
			Error: "Invalid Token",
		}, nil
	}

	if err := srv.repo.InvalidateToken(userID, tokenID); err != nil {
		log.Println("Invalidating used token error: ", err.Error())
		return &pb.TokenPair{
			Error: err.Error(),
		}, err
	}

	refreshData, err := srv.generateRefreshToken(userID)
	if err != nil {
		log.Println("generate Refresh Error: ", err.Error())
		return &pb.TokenPair{
			Error: err.Error(),
		}, err
	}

	if err := srv.repo.SaveToken(fmt.Sprintf("%s:%s:%s", userID, tokenID, refreshData.ID.String()), srv.refreshTokenDuration); err != nil {
		log.Println("saving refresh Error: ", err.Error())
		return &pb.TokenPair{
			Error: err.Error(),
		}, err
	}

	access, err := srv.generateAccessToken(userID)
	if err != nil {
		log.Println("generate Access Error: ", err.Error())
		return &pb.TokenPair{
			Error: err.Error(),
		}, err
	}
	return &pb.TokenPair{
		AccessToken:  access,
		RefreshToken: refreshData.Token,
	}, nil
}

func (srv *TokenService) DeleteUserToken(ctx context.Context, refresh *pb.RefreshToken) (*pb.Msg, error) {
	log.Printf("Receive message body from client: %s", refresh.Token)

	token, err := jwt.ParseWithClaims(refresh.GetToken(), &jwt.StandardClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(srv.refreshTokenSecret), nil
		})
	if err != nil {
		log.Println("Parsing jwt error: ", err.Error())
		return &pb.Msg{
			Error: err.Error(),
		}, err
	}
	userID := token.Claims.(*jwt.StandardClaims).Subject
	tokenID := token.Claims.(*jwt.StandardClaims).Id
	if err := srv.repo.InvalidateToken(userID, tokenID); err != nil {
		log.Println("Invalidating token error: ", err.Error())
		return &pb.Msg{
			Error: err.Error(),
		}, err
	}
	return &pb.Msg{}, nil
}
