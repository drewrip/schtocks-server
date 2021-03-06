package server

import(
	"fmt"
	"time"
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
	//"github.com/drewrip/schtocks-server/prices"
	"github.com/drewrip/schtocks-server/stocks"
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

func (s *Server) AddStockPrice(st *stocks.Stock) {
	insertPriceSQL := fmt.Sprintf(`INSERT INTO %s (time, price) values (?, ?);`, st.Ticker)


	stmt, err := s.DB.Prepare(insertPriceSQL)
	check(err)

	_, err = stmt.Exec(time.Now().UnixNano(), st.CurrentPrice)
	check(err)
	
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
	}
}
