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
			}
			if a.Currency == "EUR" {
				item := systray.AddMenuItem(fmt.Sprintf("%s %f", a.Currency, a.Balance), "")
				item.Disable()
				items = append(items, item)
			}
		}

		mSell := systray.AddMenuItem("Sell all", "")
		mBuy := systray.AddMenuItem("Buy all", "")
		mTime := systray.AddMenuItem("Time", "")
		mTime.Disable()
		eth := 0.0
		eur := 0.0

		if err != nil {
			println(err.Error())
		}
		for _, a := range accounts {
			if a.Currency == "EUR" {
				items[0].SetTitle(fmt.Sprintf("%s %f", a.Currency, a.Available))
				eur = a.Balance
			}
			if a.Currency == "ETH" {
				items[1].SetTitle(fmt.Sprintf("%s %f", a.Currency, a.Available))
				eth = a.Balance
			}
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
						eur = a.Balance
					}
					if a.Currency == "ETH" {
						items[1].SetTitle(fmt.Sprintf("%s %f", a.Currency, a.Available))
						eth = a.Balance
					}
				}
				ticker, err := client.GetTicker("ETH-EUR")
				if err != nil {
					println(err.Error())
					continue
				}
				systray.SetTitle(fmt.Sprint(ticker.Price))
				t, err :=client.GetTime()
				if err != nil {
					println(err.Error())
					continue
				}
				mTime.SetTitle(t.ISO +" "+ fmt.Sprint(ticker.Price))

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