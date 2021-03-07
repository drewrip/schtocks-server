package server

import(
	"fmt"
	"time"
	"database/sql"
	"log"
	"net/http"
	"errors"
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


func (s *Server) GetUserSummaries() map[string](map[string]int64) {
	dump := make(map[string](map[string]int64))
	getUserSummariesSQL := `SELECT * FROM market;`

	rows, err := s.DB.Query(getUserSummariesSQL)
	check(err)

	var username string
	var ticker string
	var amount int64
	var boughtPrice float64
	var boughtTime int64
	
	for rows.Next() {
		err = rows.Scan(&username, &ticker, &amount, &boughtPrice, &boughtTime)
		if _, ok := dump[username]; !ok {
			dump[username] = make(map[string]int64)
		}
		dump[username][ticker] = amount
	}

	getListOfUsersSQL := `SELECT username FROM users`

	rows, err = s.DB.Query(getListOfUsersSQL)
	check(err)

	var nonActiveUser string

	for rows.Next() {
		rows.Scan(&nonActiveUser)
		if _, ok := dump[nonActiveUser]; !ok {
			dump[nonActiveUser] = make(map[string]int64)
		}
	}

	return dump
}
func (s *Server) GetCurrentStockPrice(tr string) float64 {
	getCurrentStockPriceSQL := fmt.Sprintf(`SELECT price FROM %s ORDER BY time DESC LIMIT 1`, tr)

	rows, err := s.DB.Query(getCurrentStockPriceSQL)
	check(err)

	var price float64

	// This should only run once
	for rows.Next() {
		err = rows.Scan(&price)
		check(err)
	}

	return price
	
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


func (s *Server) AddToUserBalance(username string, money float64){
	addToUserBalanceSQL := `UPDATE users SET money = money + ? WHERE username = ?`

	stmt, err := s.DB.Prepare(addToUserBalanceSQL)
	check(err)

	_, err = stmt.Exec(money, username)
	check(err)
}

func (s *Server) AddToUserStockAmount(username string, ticker string, amount int64){

	isPresent := s.IsStockPresentForUser(username, ticker)

	addToUserBalanceSQL := ``
	if isPresent {
		addToUserBalanceSQL = `UPDATE market SET amount = amount + ? WHERE username = ? AND ticker = ?`

		stmt, err := s.DB.Prepare(addToUserBalanceSQL)
		check(err)

		_, err = stmt.Exec(amount, username, ticker)
		check(err)
	} else {
		addToUserBalanceSQL = `INSERT INTO market (username, ticker, amount) values (?, ?, ?)`

		stmt, err := s.DB.Prepare(addToUserBalanceSQL)
		check(err)

		_, err = stmt.Exec(username, ticker, amount)
		check(err)
	}

}

func (s *Server) IsStockPresentForUser(username string, ticker string) bool {
	isStockPresentForUserSQL := `SELECT COUNT(*) FROM market WHERE ticker = ? AND username = ?`

	row := s.DB.QueryRow(isStockPresentForUserSQL, ticker, username)

	var count int

	err := row.Scan(&count)
	check(err)

	return count > 0
}

func (s *Server) GetUserBalance(username string) float64 {
	getUserBalanceSQL := `SELECT money FROM users WHERE username = ?`

	rows, err := s.DB.Query(getUserBalanceSQL, username)
	check(err)

	var balance float64

	for rows.Next() {
		err = rows.Scan(&balance)
		check(err)
	}

	return balance
}

func (s *Server) GetUserStockAmount(username string, ticker string) int64 {
	getUserBalanceSQL := `SELECT amount FROM market WHERE username = ? AND ticker = ?`

	
	rows, err := s.DB.Query(getUserBalanceSQL, username, ticker)
	check(err)

	var amount int64

	for rows.Next() {
		err = rows.Scan(&amount)
		check(err)
	}

	return amount
}


func (s *Server) SellStock(username string, ticker string, amount int64){
	if s.GetUserStockAmount(username, ticker) < amount {
		check(errors.New("user trying to sell more stock than they have"))
	}

	price := s.GetCurrentStockPrice(ticker)

	s.AddToUserBalance(username, float64(amount) * price)
	s.AddToUserStockAmount(username, ticker, -1 * amount)
}

func (s *Server) BuyStock(username string, ticker string, amount int64){
	price := s.GetCurrentStockPrice(ticker)
	if s.GetUserBalance(username) < (price * float64(amount)) {
		check(errors.New("user trying to spend more money than they have"))
	}
	s.AddToUserBalance(username, -1 * float64(amount) * price)
	s.AddToUserStockAmount(username, ticker, amount)
}

func (s *Server) AddNewUser(username string, startingMoney float64){
	addNewUserSQL := `INSERT INTO users (username, money) values (?, ?)`

	stmt, err := s.DB.Prepare(addNewUserSQL)
	check(err)

	_, err = stmt.Exec(username, startingMoney)
	check(err)
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

	router.HandleFunc("/getUserSummaries", s.getUserSummariesHandler)
	router.HandleFunc("/buy", s.buyHandler)
	router.HandleFunc("/sell", s.sellHandler)
	fmt.Printf("[SERVER] Starting http server on :3432\n")
    if err := http.ListenAndServe(":3432", router); err != nil {
        log.Fatal(err)
    }
}

func (s *Server) newUserTable(){
	newUserTableSQL := `CREATE TABLE IF NOT EXISTS users (
"username" TEXT,
"money" REAL
);`

	stmt, err := s.DB.Prepare(newUserTableSQL)
	check(err)

	_, err = stmt.Exec()
	check(err)

}


func (s *Server) newMarketTable(){
	newMarketTableSQL := `CREATE TABLE IF NOT EXISTS market (
"username" TEXT,
"ticker" TEXT,
"amount" INTEGER,
"boughtPrice" REAL,
"boughtTime" INTEGER
);`
		stmt, err := s.DB.Prepare(newMarketTableSQL)
	check(err)

	_, err = stmt.Exec()
	check(err)

}


func (s *Server) startMarket(){
	s.newMarketTable()
	s.newUserTable()

	time.Sleep(time.Second)

	s.AddNewUser("Zak", 1000.0)
	s.AddNewUser("Dhruv", 1000.0)
	s.AddNewUser("Lohith", 1000.0)
	s.AddNewUser("Drew", 1000.0)
}

func (s *Server) Start() {
	go s.startMarket()
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
