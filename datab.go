package main 

import (
	"database/sql"
	"log"
)
//To add the data of user into database "users"
func Add(db *sql.DB , entry SignupJSON ){

	rollno := entry.Rollno
	name := entry.Name
	password := entry.Password
	encrypted_password := HashAndSalt(password)
	log.Println("Inserting student data")
	if(Find(db , rollno , name) == 0){
    	addDataSQL :=`INSERT INTO users (rollno , name , password) VALUES (?,?,?)`
		statement, err  := db.Prepare(addDataSQL)
		CheckError(err)
    	statement.Exec(rollno , name , encrypted_password )
    	CheckError(err)
		log.Println("Inserting student data completed")
	} else {
		log.Println("User already exists")
	}
}

// Function to find of an antry already exists
func Find(db *sql.DB , rollno int64  , name string ) int{

    getDataSQL :="SELECT * FROM users"
    rows, err := db.Query(getDataSQL)
    CheckError(err)
	defer rows.Close()
    for rows.Next() {
        var rollno_ int64 
		var name_ string 
		var pass_ string
        rows.Scan(&rollno_ ,  &name_ , &pass_ )
		//log.Println("Rollno = ", rollno_ , " Name" ,name_)
		if rollno == rollno_  {
			return 1
		}
	}
	return 0
}

func Verify( db *sql.DB , user LoginJSON) bool {
	
    getDataSQL :="SELECT * FROM users"
    rows, err := db.Query(getDataSQL)
    CheckError(err)
	defer rows.Close()
    for rows.Next() {
        var rollno_ int64 
		var name_ string 
		var pass_ string
		rows.Scan(&rollno_ ,  &name_ , &pass_ )
		if user.Rollno == rollno_ &&  ComparePasswords(pass_ , user.Password ) {
          return true
		}
		
	}
	return false
}
//To get the data of all users from database "users"
func Get(db *sql.DB ){
	getDataSQL :="SELECT * FROM users"
    rows, err := db.Query(getDataSQL)
    CheckError(err)
	defer rows.Close()
    for rows.Next() {
        var rollno int64 
		var name string 
		var pass string
        rows.Scan(&rollno ,  &name , &pass)
        log.Println("Rollno = ", rollno , " Name" ,name , "Password" , pass)
    }
}


//create a table in the sqlite database
func CreateTable(db *sql.DB){
	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
		"rollno" integer,
		"name" TEXT, 
		"password" TEXT
	);`
	stmt, err := db.Prepare(createTableSQL )
	CheckError(err)
	log.Println("Creating Table")
	stmt.Exec()
	log.Println("Table Created")
}