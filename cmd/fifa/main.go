package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/chris124567/fifa/api"
)

const (
	captchaApiKey = ""
)

type Account struct {
	Console       api.Console `json:"console"`
	VendorID      string      `json:"vendorID"`
	UserAlias     string      `json:"userAlias"`
	PersonaID     string      `json:"personaID"`
	NxMpcidCookie string      `json:"nxMpcidCookie"`
	GaCookie      string      `json:"gaCookie"`
	RemID         string      `json:"remID"`
	ProxyURL      string      `json:"proxyURL"`
}

type Player struct {
	BuyPrice  int `json:"buyPrice"`
	SellPrice int `json:"sellPrice"`
}

type Config struct {
	Accounts      []Account `json:"accounts"`
	BuyAccounts   []int     `json:"buyAccounts"`
	WatchAccounts []int     `json:"watchAccounts"`

	Players map[string]Player `json:"players"`
}

type buyPlayer struct {
	id          int64
	tradeID     string
	buyNowPrice int
	sellPrice   int
}

func solveCaptcha(client *api.Client) error {
	captcha, err := client.CaptchaData()
	if err != nil {
		return err
	}

	if len(captcha.Blob) == 0 {
		return errors.New("no captcha")
	}

	status, err := client.SubmitCaptcha(captcha.Blob)
	if err != nil {
		return err
	}

	time.Sleep(60 * time.Second)

	var token string
	for len(token) == 0 || token == "CAPCHA_NOT_READY" {
		time.Sleep(5 * time.Second)

		token, err = client.GetToken(status.Request)
		if err != nil {
			return err
		} else if token == "ERROR_CAPTCHA_UNSOLVABLE" {
			log.Printf("%s - got token: %s\n", client.VendorID(), token)
			return errors.New("unsolvable captcha")
		}
	}
	log.Printf("%s - got token: %s\n", client.VendorID(), token)

	if err := client.CaptchaValidate(token); err != nil {
		return err
	}

	return nil
}

func round100(x int) int {
	return ((x + 99) / 100) * 100
}

func main() {
	rand.Seed(1)

	file, err := os.Open("./data/config.json")
	if err != nil {
		panic(err)
	}

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		file.Close()
		panic(err)
	}
	file.Close()

	var clients []api.Client
	for i, account := range config.Accounts {
		client, err := api.NewClient(account.Console, account.VendorID, account.UserAlias, account.PersonaID, account.NxMpcidCookie, account.GaCookie, account.RemID, captchaApiKey, account.ProxyURL)
		if err != nil {
			panic(err)
		}

		log.Println("Authenticating client:", client.VendorID())
		if err := client.Authenticate(); err != nil {
			panic(err)
		}
		log.Println("Updating config.json with new RemID:", client.RemID())

		config.Accounts[i].RemID = client.RemID()

		file, err := os.Create("./data/config.json")
		if err != nil {
			panic(err)
		}
		enc := json.NewEncoder(file)
		enc.SetIndent("", "\t")
		if err := enc.Encode(config); err != nil {
			file.Close()
			panic(err)
		}
		file.Close()

		clients = append(clients, client)
	}

	var buyClients []*api.Client
	for _, i := range config.BuyAccounts {
		buyClients = append(buyClients, &clients[i])
		(&clients[i]).Tradepile()
		log.Printf("%s - %d coins\n", clients[i].VendorID(), clients[i].Coins())
	}
	var watchClients []*api.Client
	for _, i := range config.WatchAccounts {
		watchClients = append(watchClients, &clients[i])
	}

	var bought sync.Map
	playerChannel := make(chan (buyPlayer))

	go func() {
		for {
			time.Sleep(3 * time.Hour)
			for i := range clients {
				client := &clients[i]

				log.Println("Authenticating client:", client.VendorID())
				if err := client.Authenticate(); err != nil {
					panic(err)
				}
				log.Println("Updating config.json with new RemID:", client.RemID())

				config.Accounts[i].RemID = client.RemID()

				file, err := os.Create("./data/config.json")
				if err != nil {
					panic(err)
				}
				enc := json.NewEncoder(file)
				enc.SetIndent("", "\t")
				if err := enc.Encode(config); err != nil {
					file.Close()
					panic(err)
				}
				file.Close()
			}
		}
	}()

	for _, watchClient := range watchClients {
		go func(client *api.Client) {
			for {
				var target Player
				var targetID string
				// map order is randomized so this should be fine
				for id, player := range config.Players {
					target = player
					targetID = id
					break
				}
				// 2000 requests * 6 seconds/request = 3.33333333 hours
				for i := 0; i < 2000; i++ {
					if err := client.PinEventTransferMarket(); err != nil {
						panic(err)
					}

					players, err := client.TransferMarketBuyItNowByPlayer(targetID, target.BuyPrice)
					if err != nil {
						if err := solveCaptcha(client); err != nil {
							log.Printf("%s - failed to solve captcha: %v\n", client.VendorID(), err)
							time.Sleep(1 * time.Hour)
						}
					}

					for _, player := range players.AuctionInfo {
						if player.BuyNowPrice <= target.BuyPrice {
							if _, ok := bought.Load(player.TradeIDStr); ok {
								continue
							}
							bought.Store(player.TradeIDStr, struct{}{})
							playerChannel <- buyPlayer{player.ItemData.ID, player.TradeIDStr, player.BuyNowPrice, target.SellPrice}
							log.Printf("%s - %s (%d)\n", player.SellerName, player.TradeIDStr, player.BuyNowPrice)
						}
					}
					if len(players.AuctionInfo) > 0 {
						fmt.Println()
					} else {
						log.Printf("%s - nothing found\n", client.VendorID())
					}

					time.Sleep((time.Duration(1+rand.Intn(13)) * time.Second))
				}

				log.Printf("%s - finished for the day, sleeping\n", client.VendorID())
				time.Sleep(24 * time.Hour)
			}
		}(watchClient)
	}

	for player := range playerChannel {
		if len(buyClients) == 0 {
			continue
		}
		buyClient := buyClients[rand.Intn(len(buyClients))]

		// if coins := buyClient.Coins(); coins < player.buyNowPrice {
		// 	log.Printf("We only have %d coins but need %d\n", coins, player.buyNowPrice)
		// 	continue
		// }
		log.Printf("%s - buying %s: %+v", buyClient.VendorID(), player.tradeID, player)

		bid, err := buyClient.Bid(player.tradeID, player.buyNowPrice)
		if err != nil {
			if _, err := buyClient.Tradepile(); err != nil {
				// if we can't get tradepile, assume we got a captcha
				if err := solveCaptcha(buyClient); err != nil {
					log.Printf("%s - failed to solve captcha: %v, waiting 5 minutes\n", buyClient.VendorID(), err)
					time.Sleep(5 * time.Minute)
				}
			} else {
				log.Printf("%s - failed to buy %s - someone must have got it before us or we don't have enough coins\n", buyClient.VendorID(), player.tradeID)
			}
			// we were too late!
			continue
		}
		log.Printf("%s - bid on %s: %+v\n", buyClient.VendorID(), player.tradeID, bid)

		go func() {
			time.Sleep(5 * time.Second)
			list, err := buyClient.ListItem(player.id, player.tradeID)
			if err != nil {
				panic(err)
			}
			log.Printf("%s - list %s: %+v\n", buyClient.VendorID(), player.tradeID, list)

			time.Sleep(5 * time.Second)
			auction, err := buyClient.AuctionHouse(player.sellPrice-100, player.sellPrice, 86400, player.id)
			if err != nil {
				log.Printf("%s - failed to list %s - %v, waiting and trying again\n", buyClient.VendorID(), err)

				time.Sleep(5 * time.Second)
				auction, err = buyClient.AuctionHouse(player.sellPrice-100, player.sellPrice, 86400, player.id)
				if err != nil {
					panic(err)
				}
			}
			log.Printf("%s - auction on %s: %+v\n", buyClient.VendorID(), player.tradeID, auction)

			time.Sleep(5 * time.Second)
			relist, err := buyClient.Relist()
			if err != nil {
				panic(err)
			}
			log.Printf("%s - relist: %+v", buyClient.VendorID(), relist)
		}()
	}
}
