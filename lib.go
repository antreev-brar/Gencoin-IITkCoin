package main

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func CheckError(err error) {
	// catch to error.
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func HashAndSalt(pwd string) string {

	// Use GenerateFromPassword to hash & salt pwd.

	pwdbyte := []byte(pwd)
	hash, err := bcrypt.GenerateFromPassword(pwdbyte, bcrypt.MinCost)
	CheckError(err)
	log.Println("Encrypted password :", string(hash))
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func ComparePasswords(hashedPwd string, plainPwd string) bool {
	// Since we'll be Getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	bytePlainPwd := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePlainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func CountEvents(db *sql.DB, rollno int64) int64 {
	querySQL := "SELECT * FROM transactions"

	rows, err := db.Query(querySQL)
	CheckError(err)
	defer rows.Close()
	var count int64
	count = 0
	for rows.Next() {
		var type_ string
		var torollno_ int64
		var fromrollno_ int64
		var coins_ int64
		var date string
		rows.Scan(&type_, &fromrollno_, &torollno_, &coins_, &date)
		log.Println(type_, torollno_, fromrollno_, coins_, date)
		if type_ == "reward" && torollno_ == rollno {
			count++
		}

	}
	return count
}

func Taxation(entry TransactionJSON) int64 {
	batchTo := entry.ToRollno / 1000
	batchFrom := entry.FromRollno / 1000
	var taxationrate int64
	if batchTo == batchFrom {
		taxationrate = 2
	} else {
		taxationrate = 33
	}
	taxReducedCoins := (entry.Coins * (100 - taxationrate)) / 100
	return taxReducedCoins
}
