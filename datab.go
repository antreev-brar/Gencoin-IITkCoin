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
	admin := entry.Admin
	coins := 0
	encrypted_password := HashAndSalt(password)
	log.Println("Inserting student data")
	if Find(db, rollno) == 0 {
		addDataSQL := `INSERT INTO users (rollno , name , password , coins , admin) VALUES (?,?,?,? ,?)`
		statement, err := db.Prepare(addDataSQL)
		CheckError(err)
		statement.Exec(rollno, name, encrypted_password, coins, admin)
		CheckError(err)
		log.Println("Inserting student data completed")
		return true
	} else {
		log.Println("User already exists")
		return false
	}
}

func AddTransaction(db *sql.DB, transactiontype string, entry TransactionJSON) bool {

	torollno := entry.ToRollno
	fromrollno := entry.FromRollno
	coins := entry.Coins

	log.Println("Inserting transaction data")

	addDataSQL := `INSERT INTO transactions (transactiontype , fromrollno ,  torollno , coins, timestamp) VALUES (?,?,?,?,CURRENT_TIMESTAMP)`
	statement, err := db.Prepare(addDataSQL)
	CheckError(err)
	if err != nil {
		return false
	}
	statement.Exec(transactiontype, fromrollno, torollno, coins)
	CheckError(err)
	if err != nil {
		return false
	}
	log.Println("Inserting  transaction data completed")
	return true

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
		var admin_ bool
		rows.Scan(&rollno_, &name_, &pass_, &coins_, &admin_)
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
		var admin_ bool
		rows.Scan(&rollno_, &name_, &pass_, &coins_, &admin_)
		//log.Println("Rollno = ", rollno_ , " Name" ,name_)
		if rollno == rollno_ {
			return coins_
		}
	}
	return 0
}
func IsAdmin(db *sql.DB, rollno int64) bool {
	getDataSQL := "SELECT * FROM users"
	rows, err := db.Query(getDataSQL)
	CheckError(err)
	defer rows.Close()
	for rows.Next() {
		var rollno_ int64
		var name_ string
		var pass_ string
		var coins_ int64
		var admin_ bool
		rows.Scan(&rollno_, &name_, &pass_, &coins_, &admin_)
		if rollno == rollno_ {
			return admin_
		}
	}
	return false
}

//To award coins to database
// NOTE TO SELF :add functionality to have an upperbound on the number of coins to be awarded
func AddBalance(db *sql.DB, entry BalanceJSON) bool {

	if IsAdmin(db, entry.Rollno) {
		return false
	} else {
		// Create a new context, and begin a transaction
		ctx := context.Background()
		tx, err := db.BeginTx(ctx, nil)
		CheckError(err)

		// get the variables for better readability
		rollno := entry.Rollno
		coinsToAdd := entry.Coins

		log.Println("Inserting  coins in student data")
		res, err := tx.ExecContext(ctx, "UPDATE users SET coins = coins + ? WHERE rollno=? AND coins + ? <= ? ", coinsToAdd, rollno, coinsToAdd, maxcoins)
		rows_affected, _ := res.RowsAffected()

		if err != nil || rows_affected != 1 {
			// Incase we find any error in the query execution, rollback the transaction
			tx.Rollback()
			return false
		}

		//Incase of no error commit the transaction
		err = tx.Commit()
		CheckError(err)
		transactiontype := "reward"
		var newentry TransactionJSON
		newentry.ToRollno = entry.Rollno
		newentry.FromRollno = 0
		newentry.Coins = entry.Coins
		AddTransaction(db, transactiontype, newentry)
		log.Println(CountEvents(db, rollno))
		log.Println("Inserting coins to database completed")
		return true
	}
}

//To make a transaction from one user to another
func MakeTransaction(db *sql.DB, entry TransactionJSON) bool {

	fromrollno := entry.FromRollno
	torollno := entry.ToRollno
	coinsToTransfer := entry.Coins
	coinsToTransferWithTax := Taxation(entry)

	// Create a new context, and begin a transaction
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	CheckError(err)

	log.Println("Making a transaction")
	log.Println("TAX deducted amount", Taxation(entry))
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
	res, err = tx.ExecContext(ctx, "UPDATE users SET coins = coins + ? WHERE rollno=? AND coins + ? <= ? ", coinsToTransferWithTax, torollno, coinsToTransferWithTax, maxcoins)
	rows_affected, _ = res.RowsAffected()

	if err != nil || rows_affected != 1 {
		// Incase we find any error in the query execution, rollback the transaction
		tx.Rollback()
		return false
	}

	//Incase of no error Commit them
	err = tx.Commit()
	CheckError(err)
	transactiontype := "transfer"
	AddTransaction(db, transactiontype, entry)

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
		var admin_ bool
		rows.Scan(&rollno_, &name_, &pass_, &coins_, &admin_)
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
		var admin bool
		rows.Scan(&rollno, &name, &pass, &coins, &admin)
		log.Println("Rollno = ", rollno, " Name", name, "Password", pass, "coins", coins, "admin", admin)
	}
}

//create a table in the sqlite database
func CreateTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
		"rollno" integer,
		"name" TEXT, 
		"password" TEXT,
		"coins" integer ,
		"admin" bool
	);`
	stmt, err := db.Prepare(createTableSQL)
	CheckError(err)
	log.Println("Creating Table")
	stmt.Exec()
	log.Println("Table Created")
}

//create a table for transactions in the sqlite database
func CreateTableTransactions(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS transactions (
		"transactiontype" TEXT,
		"fromrollno" integer,
		"torollno" integer,
		"coins" integer,
		"timestamp" TIMESTAMP
	);`
	stmt, err := db.Prepare(createTableSQL)
	CheckError(err)
	log.Println("Creating Table")
	stmt.Exec()
	log.Println("Table Created")
}

//create a table for redeem transactions in the sqlite database
func CreateTableRedeem(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS redeem (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"status" TEXT,
		"rollno" integer,
		"coins" integer,
		"item" TEXT,
		"timestamp" TIMESTAMP
	);`
	stmt, err := db.Prepare(createTableSQL)
	CheckError(err)
	log.Println("Creating Table")
	stmt.Exec()
	log.Println("Table Created")
}

//Function to add a redeem request to the redeem table
func AddRedeem(db *sql.DB, entry RedeemJSON) bool {
	status := "Pending"
	rollno := entry.Rollno
	item := entry.Item
	coins := entry.Coins

	log.Println("Inserting redeem data")

	addDataSQL := `INSERT INTO redeem (status , rollno ,coins , item , timestamp) VALUES (?,?,?,?,CURRENT_TIMESTAMP)`
	statement, err := db.Prepare(addDataSQL)
	CheckError(err)
	if err != nil {
		return false
	}

	statement.Exec(status, rollno, coins, item)
	CheckError(err)
	if err != nil {
		return false
	}

	log.Println("Inserting redeem data completed")
	return true
}

//Function to check if the index is valid and the row on which operation is performed is Pending
//if valid return number of coins

func ValidRedeem(db *sql.DB, ind int64) RedeemJSON {
	getDataSQL := "SELECT * FROM redeem"
	rows, err := db.Query(getDataSQL)
	CheckError(err)
	defer rows.Close()
	redeem := RedeemJSON{Coins: -1, Item: "N/A", Rollno: -1}

	for rows.Next() {
		var id int64
		var status string
		var rollno int64
		var coins int64
		var item string
		var date string
		rows.Scan(&id, &status, &rollno, &coins, &item, &date)
		log.Println("Id = ", id, "Status", status, "Rollno", rollno, "coins", coins, "item", item)
		if id == ind && status == "Pending" {
			redeem.Coins = coins
			redeem.Item = item
			redeem.Rollno = rollno
			return redeem
		}
	}
	return redeem
}

//Function to update the entry in the redeem table
func UpdateRedeem(db *sql.DB, adminentry RedeemAdminJSON, entry RedeemJSON) bool {

	id := adminentry.Index
	status := adminentry.Status
	coinsToRedeem := entry.Coins
	rollno := entry.Rollno

	// Create a new context, and begin a transaction
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	CheckError(err)

	log.Println("Processing a Redeem Request")

	////////////////////////////////////////////////////
	//Update the account from which coins are being transferred
	if status {
		res, err := tx.ExecContext(ctx, "UPDATE users SET coins = coins - ? WHERE rollno=?  AND coins - ? >=0 ", coinsToRedeem, rollno, coinsToRedeem)
		CheckError(err)
		rows_affected, err := res.RowsAffected()

		if err != nil || rows_affected != 1 {
			// Incase we find any error in the query execution, rollback the transaction
			tx.Rollback()
			status = false
			// If we have to roll back, start a new context
			tx, err = db.BeginTx(ctx, nil)
			CheckError(err)
		}
	}

	//String that will be updated on redeem
	statusredeem := "Rejected"
	if status {
		statusredeem = "Approved"
	}

	////////////////////////////////////////////////////
	//Update the account to which coins are being transferred
	res, err := tx.ExecContext(ctx, "UPDATE redeem SET status = ? WHERE id = ? ", statusredeem, id)
	rows_affected, _ := res.RowsAffected()

	if err != nil || rows_affected != 1 {
		// Incase we find any error in the query execution, rollback the transaction
		tx.Rollback()
	}

	//Incase of no error Commit them
	err = tx.Commit()
	CheckError(err)

	if status {
		log.Println("Redeem completed")
		return true
	} else {
		return false
	}

}
