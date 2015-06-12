/*
PACKAGE spellAPI
rainy @ 2015-06-12 <me@rainy.im>
*/
package main

import (
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func (app *App) createSpell(c *gin.Context) {
	// time.Time.Unix(time.Now())
	token := c.MustGet("token").(*jwt.Token)

	var req SpellCrt
	err := c.Bind(&req)
	if !err {
		c.JSON(200, gin.H{"code": 400, "msg": "Not enough params"})
	} else {
		var spell Spell
		spellID := bson.NewObjectId()

		spell.ID = spellID
		spell.Owner = token.Claims["name"].(string)
		spell.Content = req.Spell
		spell.Lang = req.Lang
		spell.Status = req.Status
		spell.Timestamp = time.Time.Unix(time.Now())
		err := app.db.C("spell").Insert(spell)
		if err != nil {
			c.JSON(200, gin.H{"code": 200, "msg": "Server Error!"})
		} else {
			c.JSON(200, gin.H{"code": 200, "msg": "OK", "spell": spellID})
		}
	}

}
func (app *App) fetchSpell(c *gin.Context) {
	vistor := getUserFromToken(mySigningKey, c)
	var spellID SpellFetch
	c.Bind(&spellID)
	var spell Spell
	err := app.db.C("spell").FindId(spellID.ID).One(&spell)
	if err != nil {
		c.JSON(200, gin.H{"code": 404, "msg": "Spell Not Found!"})
		return
	}
	if spell.Status == 1 && spell.Owner != vistor {
		c.JSON(200, gin.H{"code": 403, "msg": "Private Spell!"})
	}
	c.JSON(200, gin.H{"code": 200, "spell": spell})
}
