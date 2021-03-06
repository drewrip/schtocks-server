package prices

import (
	"math/rand"
	"time"
)

type PriceModel struct {
	startPrice float64
	currentPrice float64

	volatility float64
	driftFunc func(int64) float64

	randGen *rand.Rand

	tickNum int64
}

type Price struct {
	Name string
	Ticker string
	Price float64
}

func NewPriceModel(sPrice float64, vol float64, df func(int64) float64) *PriceModel {
	currTime := time.Now().UnixNano()
    r := rand.New(rand.NewSource(currTime))
	return &PriceModel{
		startPrice: sPrice,
		currentPrice: sPrice,
		volatility: vol,
		driftFunc: df,
		randGen: r,
		startTime: currTime
	}
}


func (pm *PriceModel) SetDriftFunc(df func(int64) float64){
	pm.driftFunc = df
}

func (pm *PriceModel) NextPrice() Price {
	pm.tickNum++
	min := pm.currentPrice - pm.volatility
	max := pm.currentPrice - pm.volatility
	rNum := min + pm.randGen.Float64()*(max-min)
	return Price{
		name: pm.name,
		ticker: pm.ticker,
		price: pm.driftFunc(pm.tickNum) + rNum
	}
}




