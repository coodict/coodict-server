/*
PACKAGE structs
rainy @ 2015-06-08 <me@rainy.im>
*/
package main

type User struct {
	Name       string `json:"name" bson:"name"`
	Mail       string `bson: "mail" json:"mail"`
	Avat       string `bson: "avat" json:"avat"`
	Salt       string `bson: "salt" json:"salt"`
	Pass       string `bson: "pass" json:"pass"`
	Coins      int    `bson: "coins" json:"coins"`
	Spells     int    `bson: "spells" json:"spells"`
	Votes      int    `bson: "votes" json:"votes"`
	IsThird    bool   `bson: "isThird" json:"isThird"`
	OpenID     string `bson: "openID" json:"openID"`
	CreateDate string `bson: "createDate" json:"createDate"`
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

type SpellCrt struct {
	Lang   string `json:"mode" binding: "required"`
	Spell  string `json:"spell" binding: "required"`
	Status int8   `json:"status" binding: "required"`
}
type Spell struct {
	Name      string `bson:"name"`
	Content   string `bson:"content"`
	Lang      string `bson:"lang"`
	Len       int8   `bson:"len"`
	Owner     string `bson:"owner"`
	Status    int8   `bson:"status"`
	Timestamp int64  `bson:"timestamp"`
	IsFork    bool   `bson:"isfork"`
	From      string `bson:"from"`
	Votes     int8   `bson:"votes"`
	Share     int8   `bson:"shares"`
	Comms     int8   `bson:"comms"`
	Forks     int8   `bson:"forks"`
	Views     int32  `bson:"views"`
}
