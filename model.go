package main

import (
	"github.com/dgrijalva/jwt-go"
)

type SignupJSON struct {
	Name     string `json:"name"`
	Rollno   int64  `json:"rollno"`
	Password string `json:"password"`
	Admin    bool   `json:"admin"`
}
type LoginJSON struct {
	Rollno   int64  `json:"rollno"`
	Password string `json:"password"`
}

type CustomClaims struct {
	Rollno int64 `json:"rollno"`
	jwt.StandardClaims
}

type LoginToken struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

type GetBalanceJSON struct {
	Rollno int64 `json:"rollno"`
}

type BalanceJSON struct {
	Rollno int64 `json:"rollno"`
	Coins  int64 `json:"coins"`
}

type TransactionJSON struct {
	FromRollno int64 `json:"fromrollno"`
	ToRollno   int64 `json:"torollno"`
	Coins      int64 `json:"coins"`
}

type RedeemJSON struct {
	Rollno int64  `json:"rollno"`
	Item   string `json:"item"`
	Coins  int64  `json:"coins"`
}
type RedeemInputJSON struct {
	Item  string `json:"item"`
	Coins int64  `json:"coins"`
}
type RedeemAdminJSON struct {
	Index  int64 `json:"index"`
	Status bool  `json:"status"`
}
