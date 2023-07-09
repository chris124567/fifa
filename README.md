Adventures in reverse engineering FIFA iOS app to snipe (automatically buying undervalued players on the transfer market then reselling them) coins from early 2023.  Made 100k coins (~$3) then I lost interest because it would be expensive to scale as accounts cost a lot of money and I had another project to work on.  I do not make all the pinEvents requests that the app does so you may be eventually banned for using this.  Client is headless and does not use Selenium or Puppeteer.

# `./fifa`

Contains an unofficial FIFA API client in Go.  Able to generate the "digsig" signature to login using the installed system `node` and `./scripts/ds.js`

To make it work now, you may have to update api.appVersion,  api.apiVersion, and ClientVersion in the UtAuthRequest (client.go line 416) to whatever the FIFA companion app currently uses.  You can find that out by logging the requests the app makes in Charles Proxy after setting up TLS interception on your device.  If you want it to solve captchas for you update `captchaApiKey` in main.go to whatever your 2captcha API key.

# `./cmd/fifa`

Uses the previously mentioned FIFA API client to set up "watch accounts" and "buy accounts."  Watch accounts repeatedly search for a given player and if the price is below the specified threshold send that information to the buy accounts to buy it.  The idea was to not have accounts with the coins on them making repeated searches as that (I assume) increases ban risk.

# `./data/config.json`
Fill in the values for "userAlias", "personaID","nxMpcidCookie","gaCookie","remID" by signing into the app while logging requests in Charles or a similar MITM proxy.

`buyAccounts` and `watchAccounts` are a list of indices.  So in the configuration file currently in the repo, the 0th account would be a buy account and a watch account.  If you add more accounts you do something like:

```
    "buyAccounts": [
        5, 6
    ],
    "watchAccounts": [
        0, 1, 2, 3, 4
    ],
```

Which would make accounts 0-4 watch accounts and 5-6 buy accounts.
