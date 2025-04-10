package cabinets

type Cabinet struct {
	UUID         string `validate:"required uuid" json:"uuid"`
	Name         string `validate:"required" json:"name"`
	ClientID     string `validate:"required" json:"clientID"`
	ClientSecret string `validate:"required" json:"clientSecret"`
	CreatedAt    string `json:"createdAt"`
}

// func (c *Cabinet) ID() string {
// 	h := md5.New()
// 	io.WriteString(h, c.ClientID)
// 	io.WriteString(h, c.ClientSecret)

// 	return string(h.Sum(nil))
// }
