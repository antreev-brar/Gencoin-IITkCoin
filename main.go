package main

import (
    "database/sql"
    "log"
    _ "github.com/mattn/go-sqlite3"
)

type User struct {
    name string
    rollno int 
}
func checkError(err error) {
	// catch to error.
	if err != nil {
		log.Fatalln(err.Error())
	}
}

//To add the data of user into database "users"
func add(db *sql.DB , rollno int  , name string ){
	log.Println("Inserting student data")
    addDataSQL :=`INSERT INTO users (rollno , name) VALUES (?,?)`
	statement, err  := db.Prepare(addDataSQL)
	checkError(err)
    statement.Exec(rollno , name)
    checkError(err)
	log.Println("Inserting student data completed")
}

//To get the data of all users from database "users"
func get(db *sql.DB ){
	getDataSQL :="SELECT * FROM users"
    rows, err := db.Query(getDataSQL)
    checkError(err)
	defer rows.Close()
    for rows.Next() {
        var rollno int 
        var name string 
        rows.Scan(&rollno ,  &name)
        log.Println("Rollno = ", rollno , " Name" ,name)
    }
}


//create a table in the sqlite database
func createTable(db *sql.DB){
	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
		"rollno" integer,
		"name" TEXT
	);`
	stmt, err := db.Prepare(createTableSQL )
	checkError(err)
	log.Println("Creating Table")
	stmt.Exec()
	log.Println("Table Created")
}
func main() {
    database, _ := sql.Open("sqlite3" , "./database.db")
	defer database.Close()
	
	createTable(database)

	per1 := User{name :" antreev" , rollno : 20}
	per2 := User{name :"rodriguez" , rollno : 200}
	per3 := User{name :"alabama" , rollno : 2000}
	per4 := User{name :"drangenstein" , rollno : 20000}  
   	add(database , per1.rollno  ,per1.name)
    add(database , per2.rollno  ,per2.name)
    add(database , per3.rollno  ,per3.name)
    add(database , per4.rollno  ,per4.name)
    get(database)

}
