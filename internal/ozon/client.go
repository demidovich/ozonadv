package ozon

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

const rootUrl = "https://api-performance.ozon.ru/api"

type Client struct {
	verbose      bool
	clientId     string
	clientSecret string
	accesstoken  string
	resty        *resty.Client
}

type Config struct {
	ClientId     string
	ClientSecret string
}

func NewClient(cfg Config) *Client {
	fmt.Println("Инициализация клиента API Озон")

	c := &Client{
		resty:        resty.New(),
		clientId:     cfg.ClientId,
		clientSecret: cfg.ClientSecret,
	}

	c.resty.SetHeader("Content-Type", "application/json")
	c.resty.SetHeader("Accept", "application/json")

	return c
}

func (c *Client) SetVerbose(value bool) {
	c.verbose = value
}

func (c *Client) Get(resource string, result any) error {
	if err := c.initAccessToken(); err != nil {
		return err
	}

	url := c.url(resource)
	c.logRequest("GET", url)

	resp, err := c.resty.R().
		SetAuthToken(c.accesstoken).
		SetResult(result).
		Get(url)

	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return errors.New(
			fmt.Sprintf("Ozon Response: %s %s", resp.Status(), resp.String()),
		)
	}

	return nil
}

func (c *Client) Post(resource string, payload any, result any) error {
	if err := c.initAccessToken(); err != nil {
		return err
	}

	url := c.url(resource)
	c.logRequest("POST", url)

	resp, err := c.resty.R().
		SetAuthToken(c.accesstoken).
		SetBody(payload).
		SetResult(result).
		Post(url)

	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("Ozon Response: %s %s", resp.Status(), resp.String())
	}

	return nil
}

func (c *Client) DownloadStatistic(url string) (data []byte, err error) {
	if err = c.initAccessToken(); err != nil {
		return
	}

	resp, err := c.resty.R().
		SetAuthToken(c.accesstoken).
		Get(url)

	if err != nil {
		return data, err
	}

	if resp.StatusCode() != http.StatusOK {
		return data, fmt.Errorf("Ozon Response: %s %s", resp.Status(), resp.String())
	}

	data = resp.Body()

	return
}

func (c *Client) initAccessToken() error {
	if c.accesstoken != "" {
		return nil
	}

	url := c.url("/client/token")
	c.logRequest("POST", url)

	payload := map[string]string{
		"client_id":     c.clientId,
		"client_secret": c.clientSecret,
		"grant_type":    "client_credentials",
	}

	result := &struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}{}

	resp, err := c.resty.R().
		SetBody(payload).
		SetResult(result).
		Post(url)

	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("Ozon Access Token Response: %s %s", resp.Status(), resp.String())
	}

	c.accesstoken = result.AccessToken

	fmt.Println("Получен токен API Озон:", c.accesstoken)
	return nil
}

func (c *Client) url(resource string) string {
	return rootUrl + resource
}

func (c *Client) logRequest(method, url string) {
	if !c.verbose {
		return
	}

	fmt.Printf("Request: %s %s\n", method, url)
}
