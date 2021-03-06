package main

import(
	"github.com/drewrip/schtocks-server/stocks"
	"github.com/drewrip/schtocks-server/server"
)


func main(){
	comps := stocks.ParseFile("stocks/sample.json")
	server := server.NewServer()

	server.NewStockTable(comps[0])
}
