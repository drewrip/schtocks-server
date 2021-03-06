package main

import(
	"fmt"
	"github.com/drewrip/schtocks-server/pkg/prices"
)


func main(){
	model := prices.NewPriceModel(423.0, 10.0, func(x int64) float64 {
		return float64(x)
	})

	for i := range 100 {
		price := model.NextPrice()
		fmt.Printf("%d\t%f\n", i, price.Price)
	}
}
