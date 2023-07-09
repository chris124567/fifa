package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	http "github.com/useflyent/fhttp"
)

type Client struct {
	mu sync.Mutex

	console   Console
	vendorID  string
	sessionID string
	userAlias int
	personaID string
	remID     string

	cookies       string
	c             http.Client
	captchaClient CaptchaClient

	coins int
	count int
}

func NewClient(console Console, vendorID, userAlias, personaID, nxMpcidCookie, gaCookie, remID, captchaApiKey, proxyURL string) (Client, error) {
	parsed, err := url.Parse(proxyURL)
	if err != nil {
		return Client{}, err
	}
	userAliasInt, err := strconv.Atoi(userAlias)
	if err != nil {
		return Client{}, err
	}

	c := http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			Proxy:           http.ProxyURL(parsed),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	cc, err := NewCaptchaClient(c, captchaApiKey, userAgentWebBrowser, proxyURL)
	if err != nil {
		return Client{}, err
	}

	return Client{
		console:       console,
		vendorID:      vendorID,
		userAlias:     userAliasInt,
		personaID:     personaID,
		remID:         remID,
		cookies:       fmt.Sprintf("_nx_mpcid=%s; _ga=%s", nxMpcidCookie, gaCookie),
		c:             c,
		captchaClient: cc,
	}, nil
}

func requestClient(c http.Client, method string, url string, headers http.Header, data io.Reader) (*http.Response, error) {
	request, err := http.NewRequest(method, url, data)
	if err != nil {
		return nil, err
	}
	request.Header = headers

	response, err := c.Do(request)
	if err != nil {
		return nil, err
	}

	// sometimes we want to manually specify due accept-encoding header for
	// ordering reasons but this breaks the transparent decoding net/http so
	// handle this manually here
	switch response.Header.Get(http.CanonicalHeaderKey("Content-Encoding")) {
	case "gzip":
		response.Body, err = newCompressedReader(response.Body, newGzipReader)
	case "deflate":
		response.Body, err = newCompressedReader(response.Body, newDeflateReader)
	case "br":
		response.Body, err = newCompressedReader(response.Body, newBrotliReader)
	}
	return response, err
}

func (c *Client) request(method string, url string, headers http.Header, data io.Reader) (*http.Response, error) {
	return requestClient(c.c, method, url, headers, data)
}

func (c *Client) PinEvent(body string) (PinEventsResponse, error) {
	response, err := c.request(http.MethodPost, "https://pin-river.data.ea.com/pinEvents", http.Header{
		"Host":              {"pin-river.data.ea.com"},
		"Accept":            {"*/*"},
		"x-ea-game-id":      {"874217"},
		"x-ea-taxv":         {"1.1"},
		"x-ea-game-id-type": {"easku"},
		"Accept-Language":   {"en-us"},
		"Accept-Encoding":   {"gzip, deflate, br"},
		"Content-Type":      {"application/json"},
		"Origin":            {"file://"},
		"User-Agent":        {userAgentWebBrowser},
		"Content-Length":    {strconv.Itoa(len(body))},
		"Connection":        {"keep-alive"},
		http.HeaderOrderKey: {
			"host",
			"accept",
			"x-ea-game-id",
			"x-ea-taxv",
			"x-ea-game-id-type",
			"accept-language",
			"accept-encoding",
			"content-type",
			"origin",
			"user-agent",
			"content-length",
			"connection",
		},
	}, strings.NewReader(body))
	if err != nil {
		return PinEventsResponse{}, err
	}
	defer closeBody(response.Body)

	var result PinEventsResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return PinEventsResponse{}, err
	}
	return result, nil
}

func (c *Client) utasGet(url string, result interface{}) error {
	response, err := c.request(http.MethodGet, fifaGamePrefixURL+url, http.Header{
		"Host":            {"utas.mob.v1.fut.ea.com"},
		"Content-Type":    {"application/json"},
		"Cookie":          {c.cookies},
		"X-UT-SID":        {c.sessionID},
		"Accept":          {"*/*"},
		"User-Agent":      {userAgentWebBrowser},
		"Accept-Language": {"en-us"},
		"Accept-Encoding": {"gzip, deflate, br"},
		"Connection":      {"keep-alive"},
		http.HeaderOrderKey: []string{
			"host",
			"content-type",
			"cookie",
			"x-ut-sid",
			"accept",
			"user-agent",
			"accept-language",
			"accept-encoding",
			"connection",
		},
	}, nil)
	if err != nil {
		return err
	}
	defer closeBody(response.Body)

	return json.NewDecoder(response.Body).Decode(result)
}

func (c *Client) utasGetMarketRelated(url string, result interface{}) error {
	response, err := c.request(http.MethodGet, fifaGamePrefixURL+url, http.Header{
		"Host":            {"utas.mob.v1.fut.ea.com"},
		"Content-Type":    {"application/json"},
		"Cookie":          {c.cookies},
		"Connection":      {"keep-alive"},
		"X-UT-SID":        {c.sessionID},
		"Accept":          {"*/*"},
		"User-Agent":      {userAgentWebBrowser},
		"Accept-Language": {"en-us"},
		"Cache-Control":   {"no-cache"},
		"Accept-Encoding": {"gzip, deflate, br"},
		http.HeaderOrderKey: {
			"host",
			"content-type",
			"cookie",
			"connection",
			"x-ut-sid",
			"accept",
			"user-agent",
			"accept-language",
			"cache-control",
			"accept-encoding",
		},
	}, nil)
	if err != nil {
		return err
	}
	defer closeBody(response.Body)

	return json.NewDecoder(response.Body).Decode(result)
}

func (c *Client) utasMarketData(method string, url string, body []byte, result interface{}) (err error) {
	response, err := c.request(method, fifaGamePrefixURL+url, http.Header{
		"Host":            {"utas.mob.v1.fut.ea.com"},
		"Content-Type":    {"application/json"},
		"Origin":          {"file://"},
		"Cookie":          {c.cookies},
		"Connection":      {"keep-alive"},
		"X-UT-SID":        {c.sessionID},
		"Accept":          {"*/*"},
		"User-Agent":      {userAgentWebBrowser},
		"Content-Length":  {strconv.Itoa(len(body))},
		"Accept-Language": {"en-us"},
		"Accept-Encoding": {"gzip, deflate, br"},
		http.HeaderOrderKey: {
			"host",
			"content-type",
			"origin",
			"cookie",
			"connection",
			"x-ut-sid",
			"accept",
			"user-agent",
			"content-length",
			"accept-language",
			"accept-encoding",
		},
	}, bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer closeBody(response.Body)

	if result == nil {
		return nil
	}
	return json.NewDecoder(response.Body).Decode(result)
}

func (c *Client) CampaignActive() (result CampaignActiveResponse, err error) {
	err = c.utasGet("/scmp/campaign/active", &result)
	return
}

func (c *Client) Stadium() (result StadiumResponse, err error) {
	err = c.utasGet("/stadium", &result)
	return
}

func (c *Client) UserMassInfo() (result UserMassInfoResponse, err error) {
	err = c.utasGet("/usermassinfo", &result)
	return
}

func (c *Client) PlayStats() (result PlayStatsResponse, err error) {
	err = c.utasGet("/playStats", &result)
	return
}

func (c *Client) LiveMessage() (result LiveMessageResponse, err error) {
	err = c.utasGet("/livemessage/template?screen=companioniosfutlivemsg", &result)
	return
}

func (c *Client) SquadList() (result SquadListResponse, err error) {
	err = c.utasGet("/squad/list", &result)
	return
}

func (c *Client) Settings() (result SettingsResponse, err error) {
	err = c.utasGet("/settings", &result)
	return
}

func (c *Client) PriceLimits(tradeIDs []string) (result []PriceLimit, err error) {
	err = c.utasGet(fmt.Sprintf("/marketdata/item/pricelimits?itemIdList=%s", strings.Join(tradeIDs, ",")), &result)
	return
}

func (c *Client) Authenticate() error {
	var sid string
	var code1 ConnectAuthResponse
	{
		response, err := c.request(http.MethodGet, fmt.Sprintf("https://accounts.ea.com/connect/auth?client_id=FIFA23_COMP_APP&response_type=code&display=web2/login&locale=en_US&machineProfileKey=%s&redirect_uri=nucleus:rest&prompt=none&release_type=prod&scope=basic.identity+offline+signin+basic.entitlement+basic.persona&registration_source=315844&authentication_source=315844", c.vendorID), http.Header{
			"Host":            {"accounts.ea.com"},
			"Accept":          {"*/*"},
			"Content-Type":    {"application/json"},
			"Cookie":          {fmt.Sprintf("remid=%s; %s", c.remID, c.cookies)},
			"User-Agent":      {userAgentWebBrowser},
			"Accept-Language": {"en-us"},
			"Accept-Encoding": {"gzip, deflate, br"},
			"Connection":      {"keep-alive"},
			http.HeaderOrderKey: {
				"host",
				"accept",
				"content-type",
				"cookie",
				"user-agent",
				"accept-language",
				"accept-encoding",
				"connection",
			},
		}, nil)
		if err != nil {
			return err
		}
		defer closeBody(response.Body)

		if err := json.NewDecoder(response.Body).Decode(&code1); err != nil {
			return err
		}

		for _, cookie := range response.Cookies() {
			if cookie.Name == "sid" {
				sid = cookie.Value
			} else if cookie.Name == "remid" {
				c.remID = cookie.Value
				log.Println("RemID:", c.remID)
			}
		}
		if len(sid) == 0 {
			return errors.New("failed to find sid cookie")
		}
	}

	var token TokenResponse
	{
		response, err := c.request(http.MethodPost, fmt.Sprintf("https://accounts.ea.com/connect/token?grant_type=authorization_code&code=%s&client_id=FIFA23_COMP_APP&client_secret=ltM2J0cYMRHJyR1wABxk2lgXkSI2OwetRFO7Yd8nC0Zf9MQQB2rTmsOPaEsARBdLqCC98XTZWcynlTM1&redirect_uri=nucleus:rest&release_type=prod", code1.Code), http.Header{
			"Host":            {"accounts.ea.com"},
			"Content-Type":    {"application/x-www-form-urlencoded"},
			"Origin":          {"file://"},
			"Cookie":          {fmt.Sprintf("sid=%s; remid=%s; %s", sid, c.remID, c.cookies)},
			"Connection":      {"keep-alive"},
			"Accept":          {"*/*"},
			"User-Agent":      {userAgentWebBrowser},
			"Accept-Language": {"en-us"},
			"Accept-Encoding": {"br, gzip, deflate"},
			"Content-Length":  {"0"},
			http.HeaderOrderKey: {
				"host",
				"content-type",
				"origin",
				"cookie",
				"connection",
				"accept",
				"user-agent",
				"accept-language",
				"accept-encoding",
				"content-length",
			},
		}, nil)
		if err != nil {
			return err
		}
		defer closeBody(response.Body)

		if err := json.NewDecoder(response.Body).Decode(&token); err != nil {
			return err
		}
	}

	if len(token.AccessToken) == 0 {
		return errors.New("failed to get access token")
	}

	var code ConnectAuthResponse
	for _, shard := range []string{"shard2", "shard3", "shard5", "ut-auth"} {
		response, err := c.request(http.MethodGet, fmt.Sprintf("https://accounts.ea.com/connect/auth?client_id=FUTWEB_BK_OL_SERVER&redirect_uri=nucleus:rest&response_type=code&access_token=%s&release_type=prod&client_sequence=%s", token.AccessToken, shard), http.Header{
			"Host":            {"accounts.ea.com"},
			"Accept":          {"*/*"},
			"Content-Type":    {"application/json"},
			"Cookie":          {fmt.Sprintf("sid=%s; remid=%s; %s", sid, c.remID, c.cookies)},
			"User-Agent":      {userAgentWebBrowser},
			"Accept-Language": {"en-us"},
			"Accept-Encoding": {"gzip, deflate, br"},
			"Connection":      {"keep-alive"},
			http.HeaderOrderKey: {
				"host",
				"accept",
				"content-type",
				"cookie",
				"user-agent",
				"accept-language",
				"accept-encoding",
				"connection",
			},
		}, nil)
		if err != nil {
			return err
		}
		defer closeBody(response.Body)

		if err := json.NewDecoder(response.Body).Decode(&code); err != nil {
			return err
		}
	}

	{
		ds, err := exec.Command("node", "./scripts/ds.js", code.Code).Output()
		if err != nil {
			return err
		}

		personaID, err := strconv.ParseInt(c.personaID, 10, 64)
		if err != nil {
			return err
		}
		utAuthRequest := UtAuthRequest{
			ClientVersion:    9,
			Ds:               strings.TrimSpace(string(ds)),
			GameSku:          c.console.ID(),
			IsReadOnly:       false,
			Locale:           "en-US",
			Method:           "authcode",
			NucleusPersonaID: personaID,
			PriorityLevel:    4,
			SKU:              "FUT23IOS",
		}
		utAuthRequest.Identification.AuthCode = code.Code
		utAuthRequest.Identification.RedirectURL = "nucleus:rest"

		body, err := json.Marshal(utAuthRequest)
		if err != nil {
			return err
		}

		response, err := c.request(http.MethodPost, "https://utas.mob.v1.fut.ea.com/ut/auth", http.Header{
			"Host":                {"utas.mob.v1.fut.ea.com"},
			"X-UT-PHISHING-TOKEN": {"0"},
			"Accept":              {"*/*"},
			"Accept-Encoding":     {"br, gzip, deflate"},
			"Cache-Control":       {"no-cache"},
			"Accept-Language":     {"en-us"},
			"Content-Type":        {"application/json"},
			"Origin":              {"file://"},
			"User-Agent":          {userAgentWebBrowser},
			"Connection":          {"keep-alive"},
			"Content-Length":      {strconv.Itoa(len(body))},
			"Cookie":              {c.cookies},
			http.HeaderOrderKey: {
				"host",
				"x-ut-phishing-token",
				"accept",
				"accept-encoding",
				"cache-control",
				"accept-language",
				"content-type",
				"origin",
				"user-agent",
				"connection",
				"content-length",
				"cookie",
			},
		}, bytes.NewReader(body))
		if err != nil {
			return err
		}
		defer closeBody(response.Body)

		var utAuthResponse UtAuthResponse
		if err := json.NewDecoder(response.Body).Decode(&utAuthResponse); err != nil {
			return err
		}
		c.sessionID = utAuthResponse.Sid
	}

	return nil
}

func (c *Client) CaptchaData() (result CaptchaDataResponse, err error) {
	err = c.utasGet("/captcha/fun/data", &result)
	return
}

func (c *Client) CaptchaValidate(token string) error {
	data, err := json.Marshal(CaptchaValidateRequest{token})
	if err != nil {
		return err
	}
	return c.utasMarketData(http.MethodPost, "/captcha/fun/validate", data, nil)
}

func (c *Client) Watchlist() (result WatchlistResponse, err error) {
	err = c.utasGetMarketRelated("/watchlist", &result)
	if err == nil {
		c.mu.Lock()
		c.coins = result.Credits
		c.mu.Unlock()
	}
	return
}

func (c *Client) Tradepile() (result TradepileResponse, err error) {
	err = c.utasGetMarketRelated("/tradepile", &result)
	if err == nil {
		c.mu.Lock()
		c.coins = result.Credits
		c.mu.Unlock()
	}
	return
}

func (c *Client) TradeStatusLite(tradeIDs []string) (result TradeStatusLiteResponse, err error) {
	err = c.utasGetMarketRelated(fmt.Sprintf("/trade/status/lite?tradeIds=%s", strings.Join(tradeIDs, ",")), &result)
	return
}

func (c *Client) TransferMarketBidByPlayer(maskedDefID string, maxBid int) (result TransferMarketResponse, err error) {
	err = c.utasGetMarketRelated(fmt.Sprintf("/transfermarket?num=21&start=0&type=player&maskedDefId=%s&macr=%d", maskedDefID, maxBid), &result)
	return
}

func (c *Client) TransferMarketBuyItNowByPlayer(maskedDefID string, maxBuy int) (result TransferMarketResponse, err error) {
	err = c.utasGetMarketRelated(fmt.Sprintf("/transfermarket?num=21&start=0&type=player&maskedDefId=%s&maxb=%d", maskedDefID, maxBuy), &result)
	return
}

func (c *Client) TradeStatus(tradeIDs []string) (result TradeStatusResponse, err error) {
	err = c.utasGetMarketRelated(fmt.Sprintf("/trade/status?tradeIds=%s", strings.Join(tradeIDs, ",")), &result)
	if err == nil {
		c.mu.Lock()
		c.coins = result.Credits
		c.mu.Unlock()
	}
	return
}

func (c *Client) Bid(tradeID string, amount int) (result BidResponse, err error) {
	data, err := json.Marshal(BidRequest{amount})
	if err != nil {
		return BidResponse{}, err
	}
	err = c.utasMarketData(http.MethodPut, fmt.Sprintf("/trade/%s/bid", tradeID), data, &result)
	if err == nil {
		c.mu.Lock()
		c.coins = result.Credits
		c.mu.Unlock()
	}
	return
}

func (c *Client) ListItem(itemID int64, tradeID string) (result ItemResponse, err error) {
	data, err := json.Marshal(ItemRequest{[]ItemData{{ID: itemID, Pile: "trade", TradeID: tradeID}}})
	if err != nil {
		return ItemResponse{}, err
	}
	err = c.utasMarketData(http.MethodPut, "/item", data, &result)
	return
}

func (c *Client) AuctionHouse(startingBid int, buyNowPrice int, duration int, itemID int64) (result AuctionHouseResponse, err error) {
	data, err := json.Marshal(AuctionHouseRequest{
		BuyNowPrice: buyNowPrice,
		Duration:    duration,
		ItemData:    ItemData{ID: itemID},
		StartingBid: startingBid,
	})
	if err != nil {
		return AuctionHouseResponse{}, err
	}
	err = c.utasMarketData(http.MethodPost, "/auctionhouse", data, &result)
	return
}

func (c *Client) Relist() (result RelistResponse, err error) {
	err = c.utasMarketData(http.MethodPut, "/auctionhouse/relist", nil, &result)
	return
}

func (c *Client) PinEventTransferMarket() error {
	c.count += 1
	tsEventTime := time.Now().UTC().Format(pinEventTimeFormat)
	tsPostTime := time.Now().UTC().Format(pinEventTimeFormat)
	response, err := c.PinEvent(fmt.Sprintf(`{"custom":{"networkAccess":"W","service_plat":"%s"},"et":"client","events":[{"type":"menu","pgid":"Transfer Market Results - List View","core":{"s":%d,"pidt":"persona","pid":"%s","didm":{"idfv":"%s"},"ts_event":"%s","en":"page_view","pidm":{"nucleus":%d}}}],"gid":0,"is_sess":true,"loc":"en-US","plat":"iOS","rel":"prod","sid":"%s","taxv":"1.1","tid":"FUT23IOS","tidt":"easku","ts_post":"%s","v":"%s"}`, c.console.Name(), c.count, c.personaID, c.vendorID, tsEventTime, c.userAlias, c.sessionID, tsPostTime, appVersion))
	if err != nil {
		return err
	}
	if response.Status != "ok" {
		return fmt.Errorf("error: got status %s", response.Status)
	}

	c.count += 1
	tsEventTime = time.Now().UTC().Format(pinEventTimeFormat)
	tsPostTime = time.Now().UTC().Format(pinEventTimeFormat)
	response, err = c.PinEvent(fmt.Sprintf(`{"custom":{"networkAccess":"W","service_plat":"%s"},"et":"client","events":[{"type":"menu","pgid":"Transfer Market Results - List View","core":{"s":%d,"pidt":"persona","pid":"%s","didm":{"idfv":"%s"},"ts_event":"%s","en":"page_view","pidm":{"nucleus":%d}}}],"gid":0,"is_sess":true,"loc":"en-US","plat":"iOS","rel":"prod","sid":"%s","taxv":"1.1","tid":"FUT23IOS","tidt":"easku","ts_post":"%s","v":"%s"}`, c.console.Name(), c.count, c.personaID, c.vendorID, tsEventTime, c.userAlias, c.sessionID, tsPostTime, appVersion))
	if err != nil {
		return err
	}
	if response.Status != "ok" {
		return fmt.Errorf("error: got status %s", response.Status)
	}

	return nil
}

func (c *Client) SubmitCaptcha(blob string) (InResponse, error) {
	return c.captchaClient.FunCaptcha(funCaptchaPublicKey, funCaptchaPageURL, funCaptchaSURL, blob)
}

func (c *Client) GetToken(id string) (string, error) {
	return c.captchaClient.GetToken(id)
}

func (c *Client) Coins() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.coins
}

func (c *Client) VendorID() string {
	return c.vendorID
}

func (c *Client) RemID() string {
	return c.remID
}
