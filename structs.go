/*
PACKAGE structs
rainy @ 2015-06-08 <me@rainy.im>
*/
package main

import "gopkg.in/mgo.v2/bson"

type Tag struct {
	Value string `json: "value" bson: "value"`
	Label string `json: "label" bson: "label"`
}
type Language struct {
	Label string `json: "Label" bson: "Label"`
	Mode  string `json: "Mode" bson: "Mode"`
}
type User struct {
	Name       string `json:"name" bson:"name"`
	Mail       string `bson: "mail" json:"mail"`
	Avat       string `bson: "avat" json:"avat"`
	Salt       string `bson: "salt" json:"salt"`
	Pass       string `bson: "pass" json:"pass"`
	Tags       []Tag  `bson: "tags" json:"tags"`
	Coins      int    `bson: "coins" json:"coins"`
	Spells     int    `bson: "spells" json:"spells"`
	Votes      int    `bson: "votes" json:"votes"`
	IsThird    bool   `bson: "isThird" json:"isThird"`
	OpenID     string `bson: "openID" json:"openID"`
	CreateDate string `bson: "createDate" json:"createDate"`
}
type Spell struct {
	ID        bson.ObjectId `json:"_id" bson:"_id"`
	Name      string        `json:"name" bson:"name"`
	Content   string        `json:"content" bson:"content"`
	Lang      Language      `json:"lang" bson:"lang"`
	Tags      []Tag         `json: "tags" bson: "tags"`
	Len       int8          `json:"len" bson:"len"`
	Owner     string        `json:"owner" bson:"owner"`
	Status    int8          `json:"status" bson:"status"`
	Timestamp int64         `json:"timestamp" bson:"timestamp"`
	IsFork    bool          `json:"isfork" bson:"isfork"`
	From      string        `json:"from" bson:"from"`
	Votes     int8          `json:"votes" bson:"votes"`
	Share     int8          `json:"shares" bson:"shares"`
	Comms     int8          `json:"comms" bson:"comms"`
	Forks     int8          `json:"forks" bson:"forks"`
	Views     int32         `json:"views" bson:"views"`
}
type Promo struct {
	Code   string `json:"code" bson:"code"`
	Isused bool   `json:"isused" bson:"isused"`
}

// Form structs
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
	ID     bson.ObjectId `json:"_id"`
	Name   string        `json:"name"`
	Lang   Language      `json:"lang" binding: "required"`
	Tags   []Tag         `json: "tags" bson: "tags"`
	Spell  string        `json:"content" binding: "required"`
	Status int8          `json:"status" binding: "required"`
}

type SpellFetch struct {
	ID bson.ObjectId `json:"id" binding:"required"`
}
type SpellsOfMine struct {
	Page int8 `json:"page", binding: "required"`
	Pgsz int8 `json:"pgsz", binding: "required"`
}
