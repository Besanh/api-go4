package repository

import (
	"encoding/json"
	"go4/common/model"
	IRedis "go4/internal/redis"
)

const redisAccessTokenUser = "access_token_user"
const redisAccessTokenKey = "access_token_key"

type AuthRepository struct {
}

func NewAuthRepository() AuthRepository {
	return AuthRepository{}
}

var AuthRepo AuthRepository

/**
* Lay access token trong redis
* Parse JSON encode to v
 */
func (repo *AuthRepository) GetAccessTokenFromCache(clientId string) (interface{}, error) {
	res, err := IRedis.Redis.HMGet(redisAccessTokenUser, clientId)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	} else {
		accessTokenResponse := model.AccessToken{}
		accessToken, ok := res[0].(string)
		if ok {
			err := json.Unmarshal([]byte(accessToken), &accessTokenResponse)
			if err != nil {
				return nil, err
			}
		}
		return accessTokenResponse, nil
	}
}

/**
* Truyen token struct vao
* Insert vao redis
 */
func (repo *AuthRepository) InsertAccessTokenCache(token model.AccessToken) error {
	clientId := token.ClientID
	accessToken := token.Token
	// Return json encode of v
	jsonEncodeToken, err := json.Marshal(token)
	if err != nil {
		return err
	}
	// Tai sao phai dung interface cho clientStoreInfo
	// Vi HMSet can truyen vao interface
	jsonEncodeString := string(jsonEncodeToken)
	clientStoreInfo := map[string]interface{}{clientId: jsonEncodeString}
	accessTokenInfo := map[string]interface{}{accessToken: jsonEncodeString}
	err = IRedis.Redis.HMSet(redisAccessTokenUser, clientStoreInfo)
	if err != nil {
		return err
	}
	err = IRedis.Redis.HMSet(redisAccessTokenKey, accessTokenInfo)
	if err != nil {
		return err
	}
	return nil
}

/**
* Delete access token by token
 */
func (repo *AuthRepository) DeleteAccessToken(token model.AccessToken) error {
	clientId := token.ClientID
	accessToken := token.Token
	err := IRedis.Redis.HMDel(redisAccessTokenUser, clientId)
	if err != nil {
		return err
	}
	err = IRedis.Redis.HMDel(redisAccessTokenKey, accessToken)
	if err != nil {
		return err
	}
	return nil
}

/**
*
 */
func (repo *AuthRepository) GetAuthFromCache(token string) (interface{}, error) {
	res, err := IRedis.Redis.HMGet(redisAccessTokenKey, token)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	} else {
		accessTokenResponse := model.AccessToken{}
		accessToken, ok := res[0].(string)
		if ok {
			err := json.Unmarshal([]byte(accessToken), &accessTokenResponse)
			if err != nil {
				return nil, err
			}
		}
		return accessTokenResponse, nil
	}
}
