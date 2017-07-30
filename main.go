package main

import (
	exchange "github.com/preichenberger/go-coinbase-exchange"
	"fmt"

	"github.com/getlantern/systray"
	"time"
	"os"
	"math"
)

func main() {
	systray.Run(onReady)
}

func onReady() {
	go func() {
		systray.SetTitle("Ticker")

		secret := os.Getenv("COINBASE_SECRET")
		key := os.Getenv("COINBASE_KEY")
		passphrase := os.Getenv("COINBASE_PASSPHRASE")

		client := exchange.NewClient(secret, key, passphrase)

		accounts, err := client.GetAccounts()
		if err != nil {
			println(err.Error())
		}
		items := []*systray.MenuItem{}
		for _, a := range accounts {
			if a.Currency == "ETH" {
				item := systray.AddMenuItem(fmt.Sprintf("%s %f", a.Currency, a.Balance), "")
				item.Disable()
				items = append(items, item)
				item2 := systray.AddMenuItem(fmt.Sprintf("%s %f", a.Currency, a.Balance), "")
				item2.Disable()
				items = append(items, item2)
			}
			if a.Currency == "EUR" {
				item := systray.AddMenuItem(fmt.Sprintf("%s %f", a.Currency, a.Balance), "")
				item.Disable()
				items = append(items, item)
				item2 := systray.AddMenuItem(fmt.Sprintf("%s %f", a.Currency, a.Balance), "")
				item2.Disable()
				items = append(items, item2)
			}
			if a.Currency == "BTC" {
				item := systray.AddMenuItem(fmt.Sprintf("%s %f", a.Currency, a.Balance), "")
				item.Disable()
				items = append(items, item)
				item2 := systray.AddMenuItem(fmt.Sprintf("%s %f", a.Currency, a.Balance), "")
				item2.Disable()
				items = append(items, item2)
			}
			if a.Currency == "LTC" {
				item := systray.AddMenuItem(fmt.Sprintf("%s %f", a.Currency, a.Balance), "")
				item.Disable()
				items = append(items, item)
				item2 := systray.AddMenuItem(fmt.Sprintf("%s %f", a.Currency, a.Balance), "")
				item2.Disable()
				items = append(items, item2)
			}
		}

		mSell := systray.AddMenuItem("Sell all", "")
		mBuy := systray.AddMenuItem("Buy all", "")
		mTime := systray.AddMenuItem("Time", "")
		mTime.Disable()
		mBalanse := systray.AddMenuItem("Balance", "")
		mBalanse.Disable()
		eth := 0.0
		eur := 0.0
		btc := 0.0
		ltc := 0.0

		if err != nil {
			println(err.Error())
		}
		ticker, err := client.GetTicker("ETH-EUR")
		if err != nil {
			println(err.Error())
		}
		systray.SetTitle(fmt.Sprint(ticker.Price))

		for {

			select {
			case <-mSell.ClickedCh:
				sell(client, eth)
			case <-mBuy.ClickedCh:
				buy(client, eur)
			case <-time.After(30 * time.Second):

				accounts, err := client.GetAccounts()
				if err != nil {
					println(err.Error())
					continue
				}
				for _, a := range accounts {
					if a.Currency == "EUR" {
						items[0].SetTitle(fmt.Sprintf("%s %f", a.Currency, a.Available))
						items[1].SetTitle(fmt.Sprintf("%s %f", a.Currency, a.Hold))
						eur = a.Balance
					}
					if a.Currency == "ETH" {
						items[2].SetTitle(fmt.Sprintf("%s %f", a.Currency, a.Available))
						items[3].SetTitle(fmt.Sprintf("%s %f", a.Currency, a.Hold))
						eth = a.Balance
					}
					if a.Currency == "BTC" {
						items[4].SetTitle(fmt.Sprintf("%s %f", a.Currency, a.Available))
						items[5].SetTitle(fmt.Sprintf("%s %f", a.Currency, a.Hold))
						btc = a.Balance
					}
					if a.Currency == "LTC" {
						items[6].SetTitle(fmt.Sprintf("%s %f", a.Currency, a.Available))
						items[7].SetTitle(fmt.Sprintf("%s %f", a.Currency, a.Hold))
						ltc = a.Balance
					}
				}
				btcTicker, err := client.GetTicker("BTC-EUR")
				if err != nil {
					println(err.Error())
					continue
				}

				ltcTicker, err := client.GetTicker("LTC-EUR")
				if err != nil {
					println(err.Error())
					continue
				}

				ethTicker, err := client.GetTicker("ETH-EUR")
				if err != nil {
					println(err.Error())
					continue
				}

				systray.SetTitle(fmt.Sprint(eth*ethTicker.Price + btc*btcTicker.Price + ltc*ltcTicker.Price + eur ))
				t, err :=client.GetTime()
				if err != nil {
					println(err.Error())
					continue
				}
				mTime.SetTitle(t.ISO)
				mBalanse.SetTitle(fmt.Sprint(eth*ethTicker.Price + btc*btcTicker.Price + ltc*ltcTicker.Price + eur ))
			}

		}
	}()
}

func buy(client *exchange.Client, fund float64) {

	ticker, err := client.GetTicker("ETH-EUR")
	if err != nil {
		println(err.Error())
		return
	}

	order := exchange.Order{
		Price: Round(ticker.Bid + 0.01, .5, 2),
		Size: Round((fund-3) /ticker.Bid,.5,4),
		Side: "buy",
		ProductId: "ETH-EUR",
	}

	_, err = client.CreateOrder(&order)
	if err != nil {
		println(err.Error())
		return
	}

}

func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

func sell(client *exchange.Client, fund float64) {

	ticker, err := client.GetTicker("ETH-EUR")
	if err != nil {
		println(err.Error())
	}

	order := exchange.Order{
		Price: Round(ticker.Ask - 0.01, .5, 2),
		Size:  Round(fund, .5, 2),
		Side: "sell",
		ProductId: "ETH-EUR",
	}
	_, err = client.CreateOrder(&order)
	if err != nil {
		println(err.Error())
	}

}