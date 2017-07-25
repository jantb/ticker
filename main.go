package main

import (
	exchange "github.com/preichenberger/go-coinbase-exchange"
	"fmt"

	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"time"
	"os"
)


func main() {
	// Should be called at the very beginning of main().
	systray.Run(onReady)
}

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("Awesome App")
	systray.SetTooltip("Lantern")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
		fmt.Println("Quit now...")
	}()

	// We can manipulate the systray in other goroutines
	go func() {
		systray.SetIcon(icon.Data)
		systray.SetTitle("Ticker")

		secret := os.Getenv("COINBASE_SECRET")
		key := os.Getenv("COINBASE_KEY")
		passphrase := os.Getenv("COINBASE_PASSPHRASE")

		client := exchange.NewClient(secret, key, passphrase)

		for {
			//
			//accounts, err := client.GetAccounts()
			//if err != nil {
			//	println(err.Error())
			//}
			//
			//for _, a := range accounts {
			//	println(a.Currency)
			//	println(a.Balance)
			//}
			ticker, err := client.GetTicker("ETH-EUR")
			if err != nil {
				println(err.Error())
			}
			fmt.Println(ticker.Ask)
			systray.SetTitle(fmt.Sprint(ticker.Ask))
			time.Sleep(60 * time.Second)
		}
	}()
}