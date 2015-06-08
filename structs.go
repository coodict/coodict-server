/*
PACKAGE structs
rainy @ 2015-06-08 <me@rainy.im>
*/
package main

type User struct {
	Name    string `json:"name" bson:"name"`
	Mail    string `bson: "mail" json:"mail"`
	Avat    string `bson: "avat" json:"avat"`
	Salt    string `bson: "salt" json:"salt"`
	Pass    string `bson: "pass" json:"pass"`
	Coins   int    `bson: "coins" json:"coins"`
	Spells  int    `bson: "spells" json:"spells"`
	Votes   int    `bson: "votes" json:"votes"`
	IsThird bool   `bson: "isThird" json:"isThird"`
	OpenID  string `bson: "openID" json:"openID"`
}
type Promo struct {
	Code   string `json:"code" bson:"code"`
	Isused bool   `json:"isused" bson:"isused"`
}
type Signup struct {
	Name  string `json:"name" bson:"name" binding:"required"`
	Mail  string `json:"mail" bson:"mail" binding:"required"`
	Pass  string `json:"pass" bson:"pass" binding:"required"`
	Promo string `json:"promo" bson:"promo" binding:"required"`
}
type Signin struct {
	Name string `json:"name" binding:"required"`
	Pass string `json:"pass" binding:"required"`
}
