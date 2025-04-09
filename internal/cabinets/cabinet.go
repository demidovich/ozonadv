package cabinets

type Cabinet struct {
	Name         string `validate:"required" json:"name"`
	ClientID     string `validate:"required" json:"clientID"`
	ClientSecret string `validate:"required" json:"clientSecret"`
}
