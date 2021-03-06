package server

import(
	//"fmt"
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
	DB *sql.DB
}

func (s *Server) NewStockTable(st *stocks.Stock) {
	newStockTableSQL := `CREATE TABLE IF NOT EXISTS ? (
		"time" INTEGER,
		"price" REAL		
	  );`


	stmt, err := s.DB.Prepare(newStockTableSQL)
	check(err)

	stmt.Exec(st.Ticker)
	
}

func NewServer() *Server{
	database, err := sql.Open("sqlite3", "./schtocks.db")
	check(err)

	return &Server{
		DB: database,
		TickTime: time.Second,
	}
}
