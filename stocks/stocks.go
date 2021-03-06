package stocks

import(
	"fmt"
	"encoding/json"
	"io/ioutil"
	"math"
	"github.com/drewrip/schtocks-server/prices"
)

type Stock struct {
	name string
	ticker string
	desc string

	currentPrice float64
	ath float64
	atl float64
	
	model *prices.PriceModel
}

type StockJSON struct {
	Name string `json:"name"`
	Ticker string `json:"ticker"`
	Desc string `json:"desc"`

	StartPrice float64 `json:"startPrice"`
	Vol float64 `json:"vol"`
	DriftType string `json:"driftType"`
}

func NewStock(n string, t string, d string, sPrice float64, vol float64, df func(int64) float64) *Stock {
	return &Stock{
		name: n,
		ticker: t,
		desc: d,
		currentPrice: sPrice,
		ath: sPrice,
		atl: sPrice,
		model: prices.NewPriceModel(sPrice, vol, df),
	}
}

func ParseFile(path string) []*Stock {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("problem parsing json stocks file: %v", err))
	}
	
	stocksFromFile := []StockJSON{}

	err = json.Unmarshal(file, &stocksFromFile)
	if err != nil {
		panic(fmt.Sprintf("problem parsing json stocks file: %v", err))
	}

	newStocks := []*Stock{}
	for _, s := range stocksFromFile {

		var drift func(int64) float64
		
		if s.DriftType == "none" {
			drift = func(x int64) float64 {
				return 0
			}
		} else if s.DriftType == "up"{
			drift = func(x int64) float64 {
				return 0.005 * (float64(x) * math.Sin(float64(x)))
			}
		} else if s.DriftType == "down"{
			drift = func(x int64) float64 {
				return -0.005 * (float64(x) * math.Sin(float64(x)))
			}
		} else {
			fmt.Println("[STOCKS] Parser didn't recognize driftType. Should be 'none', 'up' or 'down' -> defaulting to 'none'")
			drift = func(x int64) float64 {
				return 0
			}
		}

		stock := &Stock{
			name: s.Name,
			ticker: s.Ticker,
			desc: s.Desc,

			currentPrice: s.StartPrice,
			ath: s.StartPrice,
			atl: s.StartPrice,
			model: prices.NewPriceModel(s.StartPrice, s.Vol, drift),
		}
		newStocks = append(newStocks, stock)
	}

	return newStocks
}

func (s *Stock) GetATH() float64 {
	return s.ath
}

func (s *Stock) GetATL() float64 {
	return s.atl
}

