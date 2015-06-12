/*
PACKAGE spellAPI
rainy @ 2015-06-12 <me@rainy.im>
*/
package main

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func (app *App) createSpell(c *gin.Context) {
	// time.Time.Unix(time.Now())
	token, _ := c.Get("token")
	tkn := token.(*jwt.Token)

	var req SpellCrt
	err := c.Bind(&req)
	if !err {
		c.JSON(200, gin.H{"code": 400, "msg": "Not enough params"})
	} else {
		var spell Spell

		spell.Owner = tkn.Claims["name"].(string)
		spell.Content = req.Spell
		spell.Lang = req.Lang
		spell.Status = req.Status
		spell.Timestamp = time.Time.Unix(time.Now())
		app.db.C("spell").Insert(spell)
		c.JSON(200, gin.H{"code": 200, "msg": "OK"})
	}

}
