package ozon

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-resty/resty/v2"
)

const rootUrl = "https://api-performance.ozon.ru/api"

type Client struct {
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
	fmt.Println("")

	c := &Client{
		resty:        resty.New(),
		clientId:     cfg.ClientId,
		clientSecret: cfg.ClientSecret,
	}

	c.resty.SetHeader("Content-Type", "application/json")
	c.resty.SetHeader("Accept", "application/json")
	c.initAccessToken()

	return c
}

func (c *Client) Get(url string, params map[string]any, result any) error {
	return nil
}

func (c *Client) Post(url string, payload map[string]string, result any) error {
	return nil
}

func (c *Client) initAccessToken() {

	url := url("/client/token")
	logRequest("POST", url)

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
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if resp.StatusCode() != http.StatusOK {
		fmt.Fprintln(os.Stderr, "Response:", resp.Status(), resp.String())
		os.Exit(1)
	}

	resp.StatusCode()
	c.accesstoken = result.AccessToken

	fmt.Println("Получен токен API Озон:", c.accesstoken)
}

func url(resource string) string {
	return rootUrl + resource
}

func logRequest(method, url string) {
	fmt.Printf("Request: %s %s\n", method, url)
}
