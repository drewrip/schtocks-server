package main

import(
	"fmt"
	"github.com/drewrip/schtocks-server/stocks"
	"github.com/drewrip/schtocks-server/server"
)


func main(){
	cks := stocks.ParseFile("stocks/sample.json")
	serv := server.NewServer()
	serv.NewStockTable(cks[0])
}
