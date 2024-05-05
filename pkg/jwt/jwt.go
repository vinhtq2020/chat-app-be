package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go-service/pkg/convert"
	"strings"
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

func decodeAccessToken(token string) (int64, *Header, *AccessTokenPayload, error) {
	var header Header
	var payload AccessTokenPayload
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return -2, nil, nil, nil
	}

	headerJSON, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return -2, nil, nil, nil
	}
	err = json.Unmarshal(headerJSON, &header)
	if err != nil {
		return -1, nil, nil, err
	}

	payloadJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return -1, nil, nil, err
	}

	err = json.Unmarshal(payloadJSON, &payload)
	if err != nil {
		return -1, nil, nil, err
	}
	return 1, &header, &payload, nil
}

func hmacsha256(input string, secretKey string) string {
	hash := hmac.New(sha256.New, []byte(secretKey))
	hash.Write([]byte(input))
	expectedMAC := hash.Sum(nil)
	signature := base64.RawURLEncoding.EncodeToString(expectedMAC)

	return signature
}

func GenerateAccessToken(userId string, secretKey string, duration time.Duration) (string, AccessTokenPayload) {
	header := Header{Algorithm: HS256, Type: JWT}
	accessPayload := AccessTokenPayload{
		UserId: userId,
		// Username:   username,
		IssuedAt:   time.Now().Unix(),
		Expiration: time.Now().Add(duration).Unix(),
	}

	accessPayloadMap := convert.ToMapOmitEmpty(accessPayload)
	accessToken := generateToken(
		header,
		accessPayloadMap, secretKey)

	return accessToken, accessPayload
}

func GereateRefreshToken(userId string, secretKey string, duration time.Duration) (string, RefreshTokenPayload) {
	header := Header{Algorithm: HS256, Type: JWT}

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

func GenerateTokens(userId string, secretKey string, accessTokenDuration time.Duration, refreshTokenDuration time.Duration) TokenData {
	accessToken, _ := GenerateAccessToken(userId, secretKey, accessTokenDuration)
	refreshToken, _ := GereateRefreshToken(userId, secretKey, refreshTokenDuration)
	return TokenData{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "bearer",
		UserId:       userId,
	}
}

func VerifyAccessToken(accessToken string, secretKey string) (int64, error) {

	parts := strings.Split(accessToken, ".")
	if len(parts) != 3 {
		return -2, nil
	}

	// decode
	res, header, payload, err := decodeAccessToken(accessToken)
	if err != nil || res <= 0 {
		return res, err
	}

	res = verifyHeader(*header)
	if res <= 0 {
		return res, nil
	}

	// compare signature system with signature request
	signatureInput := fmt.Sprintf("%s.%s", parts[0], parts[1])

	// if Algorithm is HS256
	if header.Algorithm == HS256 {
		s1 := parts[2]
		s2 := hmacsha256(signatureInput, secretKey)
		if s1 != s2 {
			return -2, nil
		}
	}

	// check os expires
	if payload.Expiration < time.Now().Unix() {
		return -3, nil
	}
	return 1, nil
}

func verifyHeader(header Header) (result int64) {
	if header.Algorithm != HS256 {
		return -2
	}
	if header.Type != JWT {
		return -2
	}
	return 1
}
