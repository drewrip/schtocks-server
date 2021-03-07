package server

import(
	"net/http"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"errors"
	"github.com/drewrip/schtocks-server/stocks"
)

type BuyOrder struct {
	Username string `json:"username"`
	Ticker string `json:"ticker"`
	Amount int64 `json:"amount"`
}

type SellOrder struct {
	Username string `json:"username"`
	Ticker string `json:"ticker"`
	Amount int64 `json:"amount"`
}

func (s *Server) getAllPricesHandler(w http.ResponseWriter, r *http.Request){
	fmt.Println("[SERVER] HTTP Request to /getAllPrices")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	priceMap := make(map[string][]TimePricePair)
	for _,c := range s.Stocks {
		priceMap[c.Ticker] = s.GetStockPrices(c)
	}
	payload, err := json.Marshal(priceMap)
	check(err)

	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}

func (s *Server) getPriceHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	b, err := ioutil.ReadAll(r.Body)
	check(err)

	var m map[string]interface{}
    err = json.Unmarshal(b, &m)
	check(err)

	ticker, ok := m["ticker"].(string)
	if !ok {
		check(errors.New("could not assert JSON key as string"))
	}
	
	priceMap := make(map[string][]TimePricePair)
	priceMap[ticker] = s.GetStockPricesByTicker(ticker)

	payload, err2 := json.Marshal(priceMap)
	check(err2)

	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}

func (s *Server) getAllStockInfoHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	allStockInfo := []stocks.StockInfo{}
	
	for _,v := range s.Stocks {
		allStockInfo = append(allStockInfo,stocks.StockInfo{
			Name: v.Name,
			Ticker: v.Ticker,
			Desc: v.Desc,
		})
	}
	
	payload, err2 := json.Marshal(allStockInfo)
	check(err2)

	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}

func (s *Server) getStockInfoHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	b, err := ioutil.ReadAll(r.Body)
	check(err)

	var m map[string]interface{}
    err = json.Unmarshal(b, &m)
	check(err)

	ticker, ok := m["ticker"].(string)
	if !ok {
		check(errors.New("could not assert JSON key as string"))
	}
	
	corrStock := s.Stocks[ticker]

	corrInfo := stocks.StockInfo{
		Name: corrStock.Name,
		Ticker: corrStock.Ticker,
		Desc: corrStock.Desc,
	}
	
	payload, err := json.Marshal(corrInfo)
	check(err)

	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}

func (s *Server) getUserSummariesHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	
	payload, err := json.Marshal(s.GetUserSummaries())
	check(err)

	w.WriteHeader(http.StatusOK)
	w.Write(payload)	
}

func (s *Server) getUserBalancesHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	
	payload, err := json.Marshal(s.GetUserBalances())
	check(err)

	w.WriteHeader(http.StatusOK)
	w.Write(payload)	
}

func (s *Server) buyHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	b, err := ioutil.ReadAll(r.Body)
	check(err)

	var order BuyOrder

    err = json.Unmarshal(b, &order)
	check(err)
	
	s.BuyStock(order.Username, order.Ticker, order.Amount)

	w.WriteHeader(http.StatusOK)
}

func (s *Server) sellHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	b, err := ioutil.ReadAll(r.Body)
	check(err)

	var order SellOrder

    err = json.Unmarshal(b, &order)
	check(err)
	
	s.SellStock(order.Username, order.Ticker, order.Amount)

	w.WriteHeader(http.StatusOK)
}
