package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go-service/pkg/convert"
	"time"
)

func generateToken(header Header, payload map[string]interface{}, secretKey string) string {

	headerJSON, _ := json.Marshal(header)
	payloadJSON, _ := json.Marshal(payload)

	headerBase64 := base64.RawURLEncoding.EncodeToString(headerJSON)
	payloadBase64 := base64.RawURLEncoding.EncodeToString(payloadJSON)

	signatureInput := fmt.Sprintf("%s.%s", headerBase64, payloadBase64)
	signatureBase64 := hmacsha256(signatureInput, secretKey)

	token := fmt.Sprintf("%s.%s.%s", headerBase64, payloadBase64, signatureBase64)
	return token
}

func hmacsha256(input string, secretKey string) string {
	hash := hmac.New(sha256.New, []byte(secretKey))
	hash.Write([]byte(input))
	expectedMAC := hash.Sum(nil)
	signature := base64.RawURLEncoding.EncodeToString(expectedMAC)

	return signature
}

func GenerateAccessToken(userId string, username string, secretKey string, duration time.Duration) (string, AccessTokenPayload) {
	header := Header{Algorithm: "HS256", Type: "JWT"}
	accessPayload := AccessTokenPayload{
		UserId:     userId,
		Username:   username,
		IssuedAt:   time.Now().Unix(),
		Expiration: time.Now().Add(duration).Unix(),
	}

	accessPayloadMap := convert.ToMapOmitEmpty(accessPayload)
	accessToken := generateToken(
		header,
		accessPayloadMap, secretKey)

	return accessToken, accessPayload
}

func GereateRefreshToken(userId string, username string, secretKey string, duration time.Duration) (string, RefreshTokenPayload) {
	header := Header{Algorithm: "HS256", Type: "JWT"}

	refreshPayload := RefreshTokenPayload{
		UserId:     userId,
		IssuedAt:   time.Now().Unix(),
		Expiration: time.Now().Add(duration).Unix(),
	}

	refreshTokenPayloadMap := convert.ToMapOmitEmpty(refreshPayload)
	refreshToken := generateToken(
		header,
		refreshTokenPayloadMap,
		secretKey)
	return refreshToken, refreshPayload
}

func GenerateTokens(userId string, username string, secretKey string, accessTokenDuration time.Duration, refreshTokenDuration time.Duration) TokenData {
	accessToken, _ := GenerateAccessToken(userId, username, secretKey, accessTokenDuration)
	refreshToken, _ := GereateRefreshToken(userId, username, secretKey, refreshTokenDuration)
	return TokenData{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "bearer",
		UserId:       userId,
	}
}
