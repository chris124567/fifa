package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"

	http "github.com/useflyent/fhttp"
)

type CaptchaClient struct {
	key       string
	userAgent string
	proxyURL  *url.URL

	c http.Client
}

type InResponse struct {
	Status  int    `json:"status"`
	Request string `json:"request"`
}

func NewCaptchaClient(c http.Client, key, userAgent, proxyURL string) (CaptchaClient, error) {
	parsed, err := url.Parse(proxyURL)
	if err != nil {
		return CaptchaClient{}, err
	}
	return CaptchaClient{
		c:         c,
		key:       key,
		userAgent: userAgent,
		proxyURL:  parsed,
	}, nil
}

func (c *CaptchaClient) FunCaptcha(publicKey, pageURL, sURL, blob string) (InResponse, error) {
	// response, err := c.c.Get(fmt.Sprintf("http://2captcha.com/in.php?key=%s&userAgent=%s&proxy=%s&method=funcaptcha&publickey=%s&pageurl=%s&surl=%s&data[blob]=%s&json=1", c.key, c.userAgent, c.proxyURL, publicKey, pageURL, sURL, blob))
	query := url.Values{}
	query.Set("key", c.key)
	query.Set("userAgent", c.userAgent)

	query.Set("method", "funcaptcha")
	query.Set("publicKey", publicKey)
	query.Set("pageurl", pageURL)
	query.Set("surl", sURL)
	query.Set("data[blob]", blob)
	query.Set("json", "1")

	if c.proxyURL != nil {
		query.Set("proxytype", strings.ToUpper(c.proxyURL.Scheme))
		query.Set("proxy", strings.TrimPrefix(c.proxyURL.String(), c.proxyURL.Scheme+"://"))
	}

	// log.Println(fmt.Sprintf("http://2captcha.com/in.php?%s", query.Encode()))
	response, err := c.c.Get(fmt.Sprintf("http://2captcha.com/in.php?%s", query.Encode()))
	if err != nil {
		return InResponse{}, err
	}
	defer closeBody(response.Body)

	var result InResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return InResponse{}, err
	}
	return result, nil
}

func (c *CaptchaClient) GetToken(id string) (string, error) {
	response, err := c.c.Get(fmt.Sprintf("http://2captcha.com/res.php?key=%s&id=%s&action=get", c.key, id))
	if err != nil {
		return "", err
	}
	defer closeBody(response.Body)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return strings.TrimPrefix(string(body), "OK|"), nil
}
