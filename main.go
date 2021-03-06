package main

import(
	"time"
	"github.com/drewrip/schtocks-server/server"
)


func main(){
	server := server.NewServer(time.Second)
	defer server.CloseDB()
	server.Start()
}
