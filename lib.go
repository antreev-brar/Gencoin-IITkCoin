package main

import (
	//"database/sql"
	"log"
	"golang.org/x/crypto/bcrypt"
)

func CheckError(err error) {
	// catch to error.
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func HashAndSalt(pwd string ) string {
    
	// Use GenerateFromPassword to hash & salt pwd.
	
	pwdbyte := []byte(pwd)
	hash, err := bcrypt.GenerateFromPassword(pwdbyte , bcrypt.MinCost)
	CheckError(err)
	log.Println("Encrypted password :" , string(hash))
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func ComparePasswords(hashedPwd string, plainPwd string) bool {
	// Since we'll be Getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	bytePlainPwd := []byte(plainPwd )
	err := bcrypt.CompareHashAndPassword(byteHash, bytePlainPwd)
	if err != nil {
		log.Println(err)
		return false
	}
	
	return true
}