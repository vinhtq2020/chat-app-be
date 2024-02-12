package model

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Provider string

var googleProdiver Provider = "google"
var microsoftProvider Provider = "microsoft"
var githubProvider Provider = "github"

func (u *Provider) Scan(v interface{}) error {
	b, ok := v.([]byte)
	if !ok {
		return errors.New("can't scan data to store in User")
	}
	return json.Unmarshal(b, &u)
}

func (u Provider) Value() string {
	return fmt.Sprintf("%v", u)
}

func GoogleProvider() *Provider {
	return &googleProdiver
}

func MicrosoftProvider() *Provider {
	return &microsoftProvider
}

func GithubProvider() *Provider {
	return &githubProvider
}
