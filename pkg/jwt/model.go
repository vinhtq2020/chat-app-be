package jwt

import "time"

type Algorithm string

func (a Algorithm) Value() string {
	return string(a)
}

type TokenType string

func (a TokenType) Value() string {
	return string(a)
}

const (
	AccessTokenDuration            = time.Hour
	RefreshTokenDuration           = time.Hour
	HS256                Algorithm = "HS256"
	JWT                  TokenType = "JWT"
)

type Header struct {
	Algorithm Algorithm `json:"algorithm"`
	Type      TokenType `json:"type"`
}

type TokenData struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	TokenType    string `json:"tokenType"`
	UserId       string `json:"userId"`
}

// Payload for AccessTokenPayload
type AccessTokenPayload struct {
	UserId     string `json:"userId"`
	Username   string `json:"username"`
	IssuedAt   int64  `json:"iat"` // the time at which the token was issued, in Unix time (seconds since January 1, 1970 UTC).
	Expiration int64  `json:"exp"` // the expiration time of the token, in Unix time (seconds since January 1, 1970 UTC).
	// Audience   string `json:"aud"`   // the intended audience of the token (e.g., the domain of the service).
	// Issuer     string `json:"iss"`   // the entity that issued the Access Token (e.g., the authentication server).
	// Scope      string `json:"scope"` // the scope of access granted by the token.
}

// Payload for RefreshTokenPayload
type RefreshTokenPayload struct {
	UserId     string `json:"userId"`
	IssuedAt   int64  `json:"iat"` // the time at which the token was issued, in Unix time (seconds since January 1, 1970 UTC).
	Expiration int64  `json:"exp"` // the expiration time of the token, in Unix time (seconds since January 1, 1970 UTC).
	Token      string `json:"token"`
}
