package server

import(
	"fmt"
	"time"
	"database/sql"
	"log"
	"net/http"
	_ "github.com/mattn/go-sqlite3"
	"github.com/drewrip/schtocks-server/stocks"
	"github.com/gorilla/mux"
)

func check(err error){
	if err != nil {
		log.Fatalf("[SERVER] %v\n", err.Error())
	}
}

type Server struct {
	TickTime time.Duration
	Ticker *time.Ticker
	DB *sql.DB
	Stocks map[string]*stocks.Stock
}

type TimePricePair struct {
	Time int64
	Price float64
}

func (s *Server) NewStockTable(st *stocks.Stock) {
	newStockTableSQL := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		"time" INTEGER,
		"price" REAL		
	  );`, st.Ticker)


	stmt, err := s.DB.Prepare(newStockTableSQL)
	check(err)

	_, err = stmt.Exec()
	check(err)
	
}

func (s *Server) GetStockPrices(st *stocks.Stock) []TimePricePair {
	getStockPricesSQL := fmt.Sprintf(`SELECT * FROM %s`, st.Ticker)

	rows, err := s.DB.Query(getStockPricesSQL)
	check(err)

	pairs := []TimePricePair{}
	
	check(err)

	var qTime int64
	var qPrice float64
	
	for rows.Next() {
		err = rows.Scan(&qTime, &qPrice)
		check(err)
		pairs = append(pairs, TimePricePair{
			Time: qTime,
			Price: qPrice,
		})
	}
	
	rows.Close()

	return pairs
}

func (s *Server) GetStockPricesByTicker(tr string) []TimePricePair {
	getStockPricesSQL := fmt.Sprintf(`SELECT * FROM %s`, tr)

	rows, err := s.DB.Query(getStockPricesSQL)
	check(err)

	pairs := []TimePricePair{}
	
	check(err)

	var qTime int64
	var qPrice float64
	
	for rows.Next() {
		err = rows.Scan(&qTime, &qPrice)
		check(err)
		pairs = append(pairs, TimePricePair{
			Time: qTime,
			Price: qPrice,
		})
	}
	
	rows.Close()

	return pairs
}

func (s *Server) AddStockPrice(st *stocks.Stock) {
	insertPriceSQL := fmt.Sprintf(`INSERT INTO %s (time, price) values (?, ?);`, st.Ticker)


	stmt, err := s.DB.Prepare(insertPriceSQL)
	check(err)

	_, err = stmt.Exec(time.Now().UnixNano(), st.CurrentPrice)
	check(err)
	
}

func (s *Server) startStocks(){

	fmt.Printf("[SERVER] Starting stock price generating loop\n")
	listOfStocks := stocks.ParseFile("./sample.json")
	
	for _, c := range listOfStocks {
		s.Stocks[c.Ticker] = c
		s.NewStockTable(c)
	}
	
	for {
		<-s.Ticker.C
		for _, v := range s.Stocks {
			v.CurrentPrice = v.Model.NextPrice()
			s.AddStockPrice(v)
		}
	}
}

func (s *Server) startRequests(){
	router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/getAllPrices", s.getAllPricesHandler)
    router.HandleFunc("/getPrice", s.getPriceHandler)
	router.HandleFunc("/getAllStockInfo", s.getAllStockInfoHandler)
	router.HandleFunc("/getStockInfo", s.getStockInfoHandler)
	fmt.Printf("[SERVER] Starting http server on :3432\n")
    if err := http.ListenAndServe(":3432", router); err != nil {
        log.Fatal(err)
    }
}

func (s *Server) Start() {
	go s.startRequests()
	s.startStocks()
}

func (s *Server) CloseDB() {
	s.DB.Close()
}

func NewServer(ticktime time.Duration) *Server{
	database, err := sql.Open("sqlite3", "./schtocks.db")
	check(err)
	
	return &Server{
		DB: database,
		TickTime: ticktime,
		Ticker: time.NewTicker(ticktime),
		Stocks: make(map[string]*stocks.Stock),
	}
}
