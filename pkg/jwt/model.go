package jwt

type Header struct {
	Algorithm string `json:"algorithm"`
	Type      string `json:"type"`
}

type TokenData struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `hson:"refreshToken"`
}
