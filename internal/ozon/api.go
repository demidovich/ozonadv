package ozon

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

const apiHost = "https://api-performance.ozon.ru"

var ErrTooManyRequests = errors.New("Ozon 429")

type api struct {
	verbose       bool
	clientId      string
	clientSecret  string
	accessToken   *accessToken
	resty         *resty.Client
	requestsCount int
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

func newApi(cfg Config, verbose bool) *api {
	a := &api{
		verbose:      verbose,
		resty:        resty.New(),
		clientId:     cfg.ClientId,
		clientSecret: cfg.ClientSecret,
	}

	a.resty.SetHeader("Content-Type", "application/json")
	a.resty.SetHeader("Accept", "application/json")

	return a
}

// Enable verbose logging
func (a *api) SetVerbose(value bool) {
	a.verbose = value
}

// HTTP Get Request
func (a *api) Get(resource string, result any) error {
	token, err := a.validAccessToken()
	if err != nil {
		return err
	}

	url := a.Url(resource)

	a.requestsCount++
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
func (a *api) GetRaw(url string) (data []byte, err error) {
	token, err := a.validAccessToken()
	if err != nil {
		return
	}

	a.requestsCount++
	a.logRequest("GET RAW", url)

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
func (a *api) Post(resource string, payload any, result any) error {
	token, err := a.validAccessToken()
	if err != nil {
		return err
	}

	url := a.Url(resource)

	a.requestsCount++
	a.logRequest("POST", url)

	resp, err := a.resty.R().
		SetAuthToken(token).
		SetBody(payload).
		SetResult(result).
		Post(url)

	if err != nil {
		return fmt.Errorf("Ozon API: POST %s %v", url, err)
	}

	if resp.StatusCode() == http.StatusTooManyRequests {
		return fmt.Errorf("Ozon Response: %s %s: %w", resp.Status(), resp.String(), ErrTooManyRequests)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("Ozon Response: %s %s", resp.Status(), resp.String())
	}

	return nil
}

func (a *api) validAccessToken() (string, error) {
	if a.accessToken != nil && a.accessToken.Valid() {
		return a.accessToken.AccessToken, nil
	}

	url := a.Url("/client/token")

	payload := map[string]string{
		"client_id":     a.clientId,
		"client_secret": a.clientSecret,
		"grant_type":    "client_credentials",
	}

	a.requestsCount++
	a.logRequest("POST", url)

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

func (a *api) RequestsCount() int {
	return a.requestsCount
}

func (a *api) Url(resource string) string {
	return fmt.Sprintf("%s/api%s", apiHost, resource)
}

func (a *api) logRequest(method, url string) {
	if !a.verbose {
		return
	}

	fmt.Printf("Ozon API: %s %s\n", method, url)
}
