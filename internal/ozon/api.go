package ozon

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
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
func (a *api) httpGet(resource string, result any) error {
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
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return errors.New(
			fmt.Sprintf("response: %s %s", resp.Status(), resp.Body()),
		)
	}

	return nil
}

// HTTP Raw Get Request
func (a *api) httpGetRaw(url string) (data []byte, err error) {
	token, err := a.validAccessToken()
	if err != nil {
		return
	}

	if !strings.HasPrefix(url, "http") {
		url = apiHost + url
	}

	a.requestsCount++
	a.logRequest("GET RAW", url)

	resp, err := a.resty.R().
		SetAuthToken(token).
		Get(url)

	if err != nil {
		return data, fmt.Errorf("GET RAW %s %v", url, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return data, fmt.Errorf("response: %s %s", resp.Status(), resp.Body())
	}

	data = resp.Body()

	return
}

// HTTP Post Request
func (a *api) httpPost(resource string, payload any, result any) error {
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
		return err
	}

	if resp.StatusCode() == http.StatusTooManyRequests {
		return fmt.Errorf("response: %s %s: %w", resp.Status(), resp.Body(), ErrTooManyRequests)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("response: %s %s", resp.Status(), resp.Body())
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
		return "", err
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("response: %s %s", resp.Status(), resp.Body())
	}

	a.accessToken.CreatedAt = time.Now()

	fmt.Println("[ozon api] получен токен доступа")

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

	fmt.Printf("[ozon api] %s %s\n", method, url)
}
