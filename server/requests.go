package server

import(
	"net/http"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"errors"
)


func (s *Server) getAllPricesHandler(w http.ResponseWriter, r *http.Request){
	fmt.Println(s.Stocks)
	fmt.Println("[SERVER] HTTP Request to /getAllPrices")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	priceMap := make(map[string][]TimePricePair)
	for _,c := range s.Stocks {
		priceMap[c.Ticker] = s.GetStockPrices(c)
	}
	payload, err := json.Marshal(priceMap)
	check(err)
	w.Write(payload)
}

func (s *Server) getPriceHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	b, err := ioutil.ReadAll(r.Body)
	check(err)

	fmt.Println(string(b))
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
	w.Write(payload)
}
