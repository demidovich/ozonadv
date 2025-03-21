package ozon

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

const apiRoot = "https://api-performance.ozon.ru/api"

type Api struct {
	verbose      bool
	clientId     string
	clientSecret string
	accessToken  *accessToken
	resty        *resty.Client
}

type accessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	CreatedAt   time.Time
}

func (a *accessToken) Valid() bool {
	lifetime := time.Now().Sub(a.CreatedAt)
	// Токен быстро протухает, возможно это какой-то баг
	// return lifetime < (time.Duration(a.ExpiresIn-500) * time.Second)
	return lifetime < (time.Duration(120) * time.Second)
}

type Config struct {
	ClientId     string
	ClientSecret string
}

func NewApi(cfg Config) *Api {
	a := &Api{
		resty:        resty.New(),
		clientId:     cfg.ClientId,
		clientSecret: cfg.ClientSecret,
	}

	a.resty.SetHeader("Content-Type", "application/json")
	a.resty.SetHeader("Accept", "application/json")

	return a
}

// Enable verbose logging
func (a *Api) SetVerbose(value bool) {
	a.verbose = value
}

// HTTP Get Request
func (a *Api) get(resource string, result any) error {
	token, err := a.validAccessToken()
	if err != nil {
		return err
	}

	url := a.url(resource)
	a.logRequest("GET", url)

	resp, err := a.resty.R().
		SetAuthToken(token).
		SetResult(result).
		Get(url)

	if err != nil {
		return fmt.Errorf("Ozon API: GET %s %v", url, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return errors.New(
			fmt.Sprintf("Ozon Response: %s %s", resp.Status(), resp.String()),
		)
	}

	return nil
}

// HTTP Raw Get Request
func (a *Api) getRaw(url string) (data []byte, err error) {
	token, err := a.validAccessToken()
	if err != nil {
		return
	}

	resp, err := a.resty.R().
		SetAuthToken(token).
		Get(url)

	if err != nil {
		return data, fmt.Errorf("Ozon API: GET RAW %s %v", url, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return data, fmt.Errorf("Ozon Response: %s %s", resp.Status(), resp.String())
	}

	data = resp.Body()

	return
}

// HTTP Post Request
func (a *Api) post(resource string, payload any, result any) error {
	token, err := a.validAccessToken()
	if err != nil {
		return err
	}

	url := a.url(resource)
	a.logRequest("POST", url)

	resp, err := a.resty.R().
		SetAuthToken(token).
		SetBody(payload).
		SetResult(result).
		Post(url)

	if err != nil {
		return fmt.Errorf("Ozon API: POST %s %v", url, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("Ozon Response: %s %s", resp.Status(), resp.String())
	}

	return nil
}

func (a *Api) validAccessToken() (string, error) {
	if a.accessToken != nil && a.accessToken.Valid() {
		return a.accessToken.AccessToken, nil
	}

	url := a.url("/client/token")
	a.logRequest("POST", url)

	payload := map[string]string{
		"client_id":     a.clientId,
		"client_secret": a.clientSecret,
		"grant_type":    "client_credentials",
	}

	resp, err := a.resty.R().
		SetBody(payload).
		SetResult(&a.accessToken).
		Post(url)

	if err != nil {
		return "", fmt.Errorf("Ozon API: POST %s %v", url, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("Ozon API: %s Response: %s %s", url, resp.Status(), resp.String())
	}

	a.accessToken.CreatedAt = time.Now()

	fmt.Println("Ozon API: получен токен доступа")

	return a.accessToken.AccessToken, nil
}

func (a *Api) url(resource string) string {
	return apiRoot + resource
}

func (a *Api) logRequest(method, url string) {
	if !a.verbose {
		return
	}

	fmt.Printf("Ozon API Request: %s %s\n", method, url)
}
