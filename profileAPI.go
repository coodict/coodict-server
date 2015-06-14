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
	app.db.C("spell").Find(bson.M{"owner": token.Claims["name"]}).Sort("-timestamp").Limit(int(req.Pgsz)).Skip(int(req.Pgsz * (req.Page - 1))).All(&spells)
	c.JSON(200, gin.H{"code": 200, "spells": spells})
}
