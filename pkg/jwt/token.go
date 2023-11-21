package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type (
	TokenOption struct {
		AccessSecret string
		AccessExpire int64
		// RefreshSecret string
		// RefreshExpire int64
		// RefreshAfter  int64
		Fields map[string]interface{}
	}
	Token struct {
		AccessToken  string `json:"access_token"`
		AccessExpire int64  `json:"access_expire"`
		// RefreshAfter  int64  `json:"refresh_after"`
		// RefreshToken  string `json:"refresh_token"`
		// RefreshExpire int64  `json:"refresh_expire"`
	}
)

// BulidTokens return a Token (which is defined by us) was built
func BuildTokens(opt TokenOption) (Token, error) {
	var token Token
	now := time.Now().Add(-time.Minute).Unix()
	// generate a accessToken
	accessToken, err := genToken(now, opt.AccessSecret, opt.Fields, opt.AccessExpire)
	if err != nil {
		return token, fmt.Errorf("create access token failed : %w ", err)
	}
	// generate a refreshToken
	// refreshToken, err := genToken(now, opt.RefreshSecret, opt.Fields, opt.RefreshExpire)
	// if err != nil {
	// 	return token, fmt.Errorf("create refresh token failed : %w ", err)
	// }

	token.AccessToken = accessToken
	token.AccessExpire = now + opt.AccessExpire
	// token.RefreshAfter = now + opt.RefreshAfter
	// token.RefreshToken = refreshToken
	// token.RefreshExpire = now + opt.RefreshExpire
	return token, nil
}

// genToken needs secretKey,payloads and token duration to generate a token
func genToken(iat int64, secretKey string, payloads map[string]interface{}, seconds int64) (string, error) {
	//make a new claims for generate token
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	for k, v := range payloads {
		claims[k] = v
	}
	// generate a token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// return signed token string
	return token.SignedString([]byte(secretKey))
}
