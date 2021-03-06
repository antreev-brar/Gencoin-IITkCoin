package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/dgrijalva/jwt-go"
	_ "github.com/mattn/go-sqlite3"
	_ "golang.org/x/crypto/bcrypt"
)

// Create the JWT key used to create the signature
var jwtKey = []byte("THE_man_has_no_face")
var maxcoins int64 = 100000

////////////////////////////////////////////////////////////////////////////////////////////////
func main() {
	database, _ := sql.Open("sqlite3", "./database.db")
	defer database.Close()
	CreateTable(database)
	CreateTableTransactions(database)
	CreateTableRedeem(database)
	Get(database)

	http.HandleFunc("/", Servepage)
	http.HandleFunc("/signup", Signup)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/getbalance", Getbalance)
	http.Handle("/addcoins", IsAuthorized(Addcoins))
	http.Handle("/redeem", IsAuthorized(Redeem))
	http.Handle("/redeemadmin", IsAuthorized(RedeemAdmin))
	http.Handle("/transaction", IsAuthorized(Transaction))
	http.Handle("/refresh", IsAuthorized(Refresh))
	http.Handle("/secretpage", IsAuthorized(Secretpage))
	log.Println("Server up at port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))

}

////////////////////////////////////////////////////////////////////////////////////////////////////
