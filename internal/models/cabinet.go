package models

import (
	"github.com/demidovich/ozonadv/pkg/utils"
	"github.com/demidovich/ozonadv/pkg/validation"
)

type Cabinet struct {
	UUID         string `validate:"required" json:"uuid"`
	Name         string `validate:"required" json:"name"`
	ClientID     string `validate:"required" json:"clientID"`
	ClientSecret string `validate:"required" json:"clientSecret"`
	CreatedAt    string `json:"createdAt"`
}

func (c *Cabinet) Validate() error {
	return validation.ValidateStruct(c)
}

func (c *Cabinet) ClientSecretMasked(length int) string {
	return utils.StringMasked(c.ClientSecret, "***", length)
}

// func (c *Cabinet) EqualTo(other *Cabinet) bool {
// 	return c.ClientID == other.ClientID && c.ClientSecret == other.ClientSecret
// }

// func (c *Cabinet) Hash() string {
// 	h := md5.New()
// 	io.WriteString(h, c.ClientID)
// 	io.WriteString(h, c.ClientSecret)

// 	return string(h.Sum(nil))
// }
