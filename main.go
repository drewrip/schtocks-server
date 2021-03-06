package main

import(
	"time"
	"github.com/drewrip/schtocks-server/stocks"
	"github.com/drewrip/schtocks-server/server"
)


func main(){
	comps := stocks.ParseFile("stocks/sample.json")
	server := server.NewServer(time.Second)
	for _, c := range comps {
		server.NewStockTable(c)
	}

	for i := 0; i < 60; i++{
		<-server.Ticker.C
		for _, c := range comps {
			c.CurrentPrice = c.Model.NextPrice()
			server.AddStockPrice(c)
		}
		
	}
	server.CloseDB()
}
