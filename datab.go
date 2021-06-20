package main

import (
	"context"
	"database/sql"
	"log"
)

//To add the data of user into database "users"
func Add(db *sql.DB, entry SignupJSON) bool {

	rollno := entry.Rollno
	name := entry.Name
	password := entry.Password
	coins := 0
	encrypted_password := HashAndSalt(password)
	log.Println("Inserting student data")
	if Find(db, rollno) == 0 {
		addDataSQL := `INSERT INTO users (rollno , name , password , coins) VALUES (?,?,?,?)`
		statement, err := db.Prepare(addDataSQL)
		CheckError(err)
		statement.Exec(rollno, name, encrypted_password, coins)
		CheckError(err)
		log.Println("Inserting student data completed")
		return true
	} else {
		log.Println("User already exists")
		return false
	}
}

// Function to find of an antry already exists
func Find(db *sql.DB, rollno int64) int {

	getDataSQL := "SELECT * FROM users"
	rows, err := db.Query(getDataSQL)
	CheckError(err)
	defer rows.Close()
	for rows.Next() {
		var rollno_ int64
		var name_ string
		var pass_ string
		var coins_ int64
		rows.Scan(&rollno_, &name_, &pass_, &coins_)
		//log.Println("Rollno = ", rollno_ , " Name" ,name_)
		if rollno == rollno_ {
			return 1
		}
	}
	return 0
}

//Function to get the balance
func FindBalance(db *sql.DB, rollno int64) int64 {

	getDataSQL := "SELECT * FROM users"
	rows, err := db.Query(getDataSQL)
	CheckError(err)
	defer rows.Close()
	for rows.Next() {
		var rollno_ int64
		var name_ string
		var pass_ string
		var coins_ int64
		rows.Scan(&rollno_, &name_, &pass_, &coins_)
		//log.Println("Rollno = ", rollno_ , " Name" ,name_)
		if rollno == rollno_ {
			return coins_
		}
	}
	return 0
}

//To award coins to database
// NOTE TO SELF :add functionality to have an upperbound on the number of coins to be awarded
func AddBalance(db *sql.DB, entry BalanceJSON) bool {

	// Create a new context, and begin a transaction
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	CheckError(err)

	// get the variables for better readability
	rollno := entry.Rollno
	coinsToAdd := entry.Coins

	log.Println("Inserting  coins in student data")
	res, err := tx.ExecContext(ctx, "UPDATE users SET coins = coins + ? WHERE rollno=? ", coinsToAdd, rollno)
	rows_affected, _ := res.RowsAffected()

	if err != nil || rows_affected != 1 {
		// Incase we find any error in the query execution, rollback the transaction
		tx.Rollback()
		return false
	}

	//Incase of no error commit the transaction
	err = tx.Commit()
	CheckError(err)

	log.Println("Inserting coins to database completed")
	return true
}

//To make a transaction from one user to another
func MakeTransaction(db *sql.DB, entry TransactionJSON) bool {

	fromrollno := entry.FromRollno
	torollno := entry.ToRollno
	coinsToTransfer := entry.Coins

	// Create a new context, and begin a transaction
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	CheckError(err)

	log.Println("Making a transaction")

	////////////////////////////////////////////////////
	//Update the account from which coins are being transferred
	res, err := tx.ExecContext(ctx, "UPDATE users SET coins = coins - ? WHERE rollno=?  AND coins - ? >=0 ", coinsToTransfer, fromrollno, coinsToTransfer)
	CheckError(err)
	rows_affected, err := res.RowsAffected()

	if err != nil || rows_affected != 1 {
		// Incase we find any error in the query execution, rollback the transaction
		tx.Rollback()
		return false
	}

	//time.Sleep(10 * time.Second)
	////////////////////////////////////////////////////
	//Update the account to which coins are being transferred
	res, err = tx.ExecContext(ctx, "UPDATE users SET coins = coins + ? WHERE rollno=? ", coinsToTransfer, torollno)
	rows_affected, _ = res.RowsAffected()

	if err != nil || rows_affected != 1 {
		// Incase we find any error in the query execution, rollback the transaction
		tx.Rollback()
		return false
	}

	//Incase of no error Commit them
	err = tx.Commit()
	CheckError(err)

	log.Println("Transaction completed")
	return true
}

//Verify if rollno and password match in the database
func Verify(db *sql.DB, user LoginJSON) bool {

	getDataSQL := "SELECT * FROM users"
	rows, err := db.Query(getDataSQL)
	CheckError(err)
	defer rows.Close()
	for rows.Next() {
		var rollno_ int64
		var name_ string
		var pass_ string
		var coins_ int64
		rows.Scan(&rollno_, &name_, &pass_, &coins_)
		if user.Rollno == rollno_ && ComparePasswords(pass_, user.Password) {
			return true
		}

	}
	return false
}

//To get the data of all users from database "users"
func Get(db *sql.DB) {
	getDataSQL := "SELECT * FROM users"
	rows, err := db.Query(getDataSQL)
	CheckError(err)
	defer rows.Close()
	for rows.Next() {
		var rollno int64
		var name string
		var pass string
		var coins int64
		rows.Scan(&rollno, &name, &pass, &coins)
		log.Println("Rollno = ", rollno, " Name", name, "Password", pass, "coins", coins)
	}
}

//create a table in the sqlite database
func CreateTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
		"rollno" integer,
		"name" TEXT, 
		"password" TEXT,
		"coins" integer
	);`
	stmt, err := db.Prepare(createTableSQL)
	CheckError(err)
	log.Println("Creating Table")
	stmt.Exec()
	log.Println("Table Created")
}
