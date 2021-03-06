package main

import(
	"fmt"
	"github.com/drewrip/schtocks-server/stocks"
)


func main(){
	s := stocks.ParseFile("stocks/sample.json")
	fmt.Println(s)
}
