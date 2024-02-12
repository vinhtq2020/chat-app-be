package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
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

func hmacsha256(input string, secretKey string) []byte {
	hash := hmac.New(sha256.New, []byte(secretKey))
	hash.Write([]byte(input))
	return hash.Sum(nil)
}

func GenerateTokens(userId string, username string, secretKey string) TokenData {
	accessToken := generateToken(
		Header{Algorithm: "HS256", Type: "JWT"},
		map[string]interface{}{
			"userId":   userId,
			"username": username,
			"exp":      time.Minute * 15,
		}, secretKey)
	refreshToken := generateToken(
		Header{Algorithm: "HS256", Type: "JWT"},
		map[string]interface{}{
			"userId": userId,
			"exp":    time.Hour * 24 * 7,
		}, secretKey)
	return TokenData{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
