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
		c.JSON(200, gin.H{"code": 400, "msg": "参数错误！"})
	} else {
		var spell Spell
		if req.ID != "" {
			app.db.C("spell").FindId(req.ID).One(&spell)
			if spell.ID != "" {
				change := bson.M{"$set": bson.M{"name": req.Name, "content": req.Spell, "mode": req.Lang, "label": req.Label, "status": req.Status, "timestamp": time.Time.Unix(time.Now())}}
				err := app.db.C("spell").UpdateId(req.ID, change)
				if err != nil {
					c.JSON(200, gin.H{"code": 500, "msg": "服务器错误，请稍后再试！"})
					return
				}
				c.JSON(200, gin.H{"code": 200, "msg": "OK", "spell": spell.ID})
				return
			}
		}
		spellID := bson.NewObjectId()

		spell.ID = spellID
		spell.Name = req.Name
		spell.Owner = token.Claims["name"].(string)
		spell.Content = req.Spell
		spell.Lang = req.Lang
		spell.Label = req.Label
		spell.Status = req.Status
		spell.Timestamp = time.Time.Unix(time.Now())
		err := app.db.C("spell").Insert(spell)
		if err != nil {
			c.JSON(200, gin.H{"code": 200, "msg": "服务器错误，请稍后再试！"})
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
		c.JSON(200, gin.H{"code": 404, "msg": "404"})
		return
	}
	if spell.Status == 1 && spell.Owner != vistor {
		c.JSON(200, gin.H{"code": 403, "msg": "主人未公开"})
	}
	c.JSON(200, gin.H{"code": 200, "spell": spell})
}
