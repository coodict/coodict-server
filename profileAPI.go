/*
PACKAGE profileAPI
rainy @ 2015-06-13 <me@rainy.im>
*/
package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

func (app *App) mySpells(c *gin.Context) {
	token := c.MustGet("token").(*jwt.Token)

	var req SpellsOfMine
	c.Bind(&req)
	spells := make([]Spell, req.Pgsz)
	query := bson.M{"owner": token.Claims["name"], "status": bson.M{"$lt": STATUS_DELETE}}
	if req.Lang != "" {
		query["lang.mode"] = req.Lang
	}
	if req.Tag != "" {
		query["tags.value"] = req.Tag
	}
	var sorts string
	switch req.Sorts {
	case "late":
		sorts = "-timestamp"
	case "earl":
		sorts = "timestamp"
	case "view":
		sorts = "-views"
	}

	app.db.C("spell").Find(query).Sort(sorts).Limit(int(req.Pgsz)).Skip(int(req.Pgsz * (req.Page - 1))).All(&spells)
	c.JSON(200, gin.H{"code": 200, "spells": spells})
}
