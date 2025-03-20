package ozon

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

const apiRoot = "https://api-performance.ozon.ru/api"

type Client struct {
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

func NewClient(cfg Config) *Client {
	c := &Client{
		resty:        resty.New(),
		clientId:     cfg.ClientId,
		clientSecret: cfg.ClientSecret,
	}

	c.resty.SetHeader("Content-Type", "application/json")
	c.resty.SetHeader("Accept", "application/json")

	return c
}

// Enable verbose logging
func (c *Client) SetVerbose(value bool) {
	c.verbose = value
}

// HTTP Get Request
func (c *Client) get(resource string, result any) error {
	token, err := c.validAccessToken()
	if err != nil {
		return err
	}

	url := c.url(resource)
	c.logRequest("GET", url)

	resp, err := c.resty.R().
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
func (c *Client) getRaw(url string) (data []byte, err error) {
	token, err := c.validAccessToken()
	if err != nil {
		return
	}

	resp, err := c.resty.R().
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
func (c *Client) post(resource string, payload any, result any) error {
	token, err := c.validAccessToken()
	if err != nil {
		return err
	}

	url := c.url(resource)
	c.logRequest("POST", url)

	resp, err := c.resty.R().
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

func (c *Client) validAccessToken() (string, error) {
	if c.accessToken != nil && c.accessToken.Valid() {
		return c.accessToken.AccessToken, nil
	}

	url := c.url("/client/token")
	c.logRequest("POST", url)

	payload := map[string]string{
		"client_id":     c.clientId,
		"client_secret": c.clientSecret,
		"grant_type":    "client_credentials",
	}

	resp, err := c.resty.R().
		SetBody(payload).
		SetResult(&c.accessToken).
		Post(url)

	if err != nil {
		return "", fmt.Errorf("Ozon API: POST %s %v", url, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("Ozon API: %s Response: %s %s", url, resp.Status(), resp.String())
	}

	c.accessToken.CreatedAt = time.Now()

	fmt.Println("Ozon API: получен токен доступа")

	return c.accessToken.AccessToken, nil
}

func (c *Client) url(resource string) string {
	return apiRoot + resource
}

func (c *Client) logRequest(method, url string) {
	if !c.verbose {
		return
	}

	fmt.Printf("Ozon API Request: %s %s\n", method, url)
}
