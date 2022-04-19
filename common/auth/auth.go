package auth

import (
	"encoding/base64"
	"errors"
	"strings"
	"time"

	"go4/common/log"
	"go4/common/model"
	"go4/repository"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type Config struct {
	ExpiredTime int
	TokenType   string
}

var AuthConfig Config

type AuthClient struct {
	ClientID     string
	ClientSecret string
	UserId       string
	Scope        string
	User         model.User
}

type AccessTokenResponse struct {
	UserID       string `json:"user_id"`
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiredTime  int    `json:"expired_at"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
}

func NewAuthUtil(config Config) {
	AuthConfig.ExpiredTime = config.ExpiredTime
	AuthConfig.TokenType = config.TokenType
}

func ClientCredential(client AuthClient, isRefresh bool) (interface{}, error) {
	accessToken, err := CheckAccessTokenCache(client)
	if err != nil {
		log.Error("Auth Middleware", "ClientCredential", err)
		return nil, err
	}
	accessTokenResponse, err := CreateResponseAccessToken(accessToken, isRefresh)
	if err != nil {
		log.Error("Auth Middleware", "ClientCredential - CreateResponseAccessToken", err)
		return nil, err
	}
	return accessTokenResponse, nil
}

func CreateResponseAccessToken(token model.AccessToken, isRefresh bool) (AccessTokenResponse, error) {
	response := AccessTokenResponse{}
	if token.Token == "" {
		return response, errors.New("token is null")
	}
	if !isRefresh {
		response = AccessTokenResponse{
			UserID:       token.UserID,
			Token:        token.Token,
			RefreshToken: token.RefreshToken,
			TokenType:    AuthConfig.TokenType,
			Scope:        token.Scope,
		}
	} else {
		response = AccessTokenResponse{
			UserID:       token.UserID,
			Token:        token.Token,
			RefreshToken: token.RefreshToken,
			ExpiredTime:  token.ExpiredTime,
			TokenType:    AuthConfig.TokenType,
			Scope:        token.Scope,
		}
	}

	return response, nil
}

func CheckAccessTokenCache(client AuthClient) (model.AccessToken, error) {
	var accessToken model.AccessToken
	log.Info("Auth Middleware", "CheckAccessTokenCache - clientID", client.ClientID)
	accessTokenRes, err := repository.AuthRepo.GetAccessTokenFromCache(client.ClientID)
	log.Info("Auth Middleware", "accessTokenRes", accessTokenRes)
	if err != nil {
		log.Error("Auth Middleware", "CheckAccessTokenCache - AuthRepo GetAccessTokenFromCache", err)
		return accessToken, err
	}
	if accessTokenRes != "" {
		// sos
		accessToken, ok := accessTokenRes.(model.AccessToken)
		if !ok {
			log.Error("Auth Middleware", "CheckAccessTokenCache", err)
			return accessToken, err
		}
	}

	if accessToken.ClientID == "" {
		accessToken = CreateAccessToken(client)
		log.Info("Auth Middleware", "CheckAccessToken - CreateAccessToken", accessToken)
		err := repository.AuthRepo.InsertAccessTokenCache(accessToken)
		if err != nil {
			log.Error("Auth Middleware", "CheckAccessTokenCache - InsertAccessTokenCache", err)
			return accessToken, err
		}
	} else {
		timeIn := accessToken.CreatedAt.Add(time.Second * time.Duration(accessToken.ExpiredTime))
		if timeIn.Sub(time.Now().Local()) <= 0 {
			log.Info("Auth Middleware", "CheckAccessTokenCache", "accesstoken already expired and create new accesstoken")
			err := repository.AuthRepo.DeleteAccessToken(accessToken)
			if err != nil {
				return accessToken, err
			}
			accessToken = CreateAccessToken(client)
			log.Info("Auth Middleware", "CheckAccessTokenCache - CreateAccessToken 2", accessToken)
			err = repository.AuthRepo.InsertAccessTokenCache(accessToken)
			if err != nil {
				log.Error("Auth Middleware", "ChecAccessTokenCache - InsertAccessTokenCache", err)
				return accessToken, err
			}
		} else {
			log.Info("Auth Middleware", "CheckAccessTokenCache", "accesstoken already existed")
		}
	}
	return accessToken, nil
}

func CreateAccessToken(client AuthClient) model.AccessToken {
	jwtData := make(map[string]string)
	jwtData["username"] = client.User.Username
	jwtData["level"] = client.User.Level

	accessToken := model.AccessToken{
		ClientID:     client.ClientID,
		UserID:       client.UserId,
		Token:        GenerateToken(client.ClientID),
		RefreshToken: GenerateRefreshToken(client.ClientID),
		CreatedAt:    time.Now().Local(),
		ExpiredTime:  AuthConfig.ExpiredTime,
		Scope:        client.Scope,
		TokenType:    AuthConfig.TokenType,
		JWT:          GenerateJWT(client.ClientID, jwtData),
	}

	return accessToken
}

// Tao token
func GenerateToken(id string) string {
	uuidNew, _ := uuid.NewRandom()
	idEnc := base64.StdEncoding.EncodeToString([]byte(id))
	token := strings.Replace(uuidNew.String(), "-", "", -1)
	token = token + "-" + idEnc
	return token
}

// Tao refresh token
func GenerateRefreshToken(id string) string {
	return GenerateToken(id)
}

// Tao jwt
func GenerateJWT(id string, data map[string]string) string {
	claim := jwt.MapClaims{
		"iss":  "api",
		"sub":  "anhle",
		"auth": "anhle",
		"jti":  id,
		"id":   id,
	}
	if len(data) > 0 {
		for key, value := range data {
			claim[key] = value
		}
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claim)
	jwtToken, _ := token.SignedString([]byte("secret"))

	return jwtToken
}
