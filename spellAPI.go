/*
PACKAGE spellAPI
rainy @ 2015-06-12 <me@rainy.im>
*/
package main

import (
	"fmt"
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
	fmt.Printf("%+v", req)
	if err != nil {
		c.JSON(200, gin.H{"code": 400, "msg": "参数错误！"})
	} else {
		var spell Spell
		// Update user metas
		if len(req.Tags) > 0 {
			fmt.Println(req.Tags)
			app.db.C("user").Update(bson.M{"name": token.Claims["name"]}, bson.M{"$addToSet": bson.M{"tags": bson.M{"$each": req.Tags}}})
		}
		app.db.C("user").Update(bson.M{"name": token.Claims["name"]}, bson.M{"$addToSet": bson.M{"langs": req.Lang}})

		if req.ID != "" {
			// Update api
			app.db.C("spell").FindId(req.ID).One(&spell)
			if spell.ID != "" {
				change := bson.M{"$set": bson.M{"tags": req.Tags, "name": req.Name, "content": req.Spell, "lang": req.Lang, "status": req.Status, "timestamp": time.Time.Unix(time.Now())}}
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
		fmt.Println(req.Tags)

		spell.ID = spellID
		spell.Name = req.Name
		spell.Owner = token.Claims["name"].(string)
		spell.Content = req.Spell
		spell.Lang = req.Lang
		spell.Status = req.Status
		spell.Tags = req.Tags
		spell.Timestamp = time.Time.Unix(time.Now())
		err := app.db.C("spell").Insert(spell)
		if err != nil {
			c.JSON(200, gin.H{"code": 500, "msg": "服务器错误，请稍后再试！"})
		} else {
			c.JSON(200, gin.H{"code": 200, "msg": "OK", "spell": spellID})
		}
	}

}
func (app *App) square(c *gin.Context) {
	var hots, news []Spell

	query := bson.M{"status": STATUS_PUBLIC}

	app.db.C("spell").Find(query).Sort("-views").Limit(10).All(&hots)
	app.db.C("spell").Find(query).Sort("-timestamp").Limit(10).All(&news)
	c.JSON(200, gin.H{"code": 200, "hots": hots, "news": news})
}
func (app *App) fetchSpell(c *gin.Context) {
	vistor := getUserFromToken(mySigningKey, c)
	var req SpellFetch
	c.Bind(&req)
	var spell Spell
	err := app.db.C("spell").Find(bson.M{"_id": req.ID, "status": bson.M{"$ne": STATUS_DELETE}}).One(&spell)
	if err != nil {
		c.JSON(200, gin.H{"code": 404, "msg": "404"})
		return
	}
	if spell.Status == STATUS_PRIVATE && spell.Owner != vistor {
		c.JSON(200, gin.H{"code": 403, "msg": "主人未公开"})
		return
	}
	app.db.C("spell").UpdateId(req.ID, bson.M{"$inc": bson.M{"views": 1}})
	c.JSON(200, gin.H{"code": 200, "spell": spell})
}

func (app *App) deleteSpell(c *gin.Context) {
	token := c.MustGet("token").(*jwt.Token)

	var req SpellFetch
	err := c.Bind(&req)
	if err != nil {
		c.JSON(200, gin.H{"code": 400, "msg": "参数错误！"})
		return
	}
	var spell Spell
	err2 := app.db.C("spell").FindId(req.ID).One(&spell)
	if err2 != nil {
		c.JSON(200, gin.H{"code": 404, "msg": "404"})
		return
	}
	if spell.Owner != token.Claims["name"] {
		c.JSON(200, gin.H{"code": 403, "msg": "不是你的你别管！"})
		return
	}
	app.db.C("spell").UpdateId(req.ID, bson.M{"$set": bson.M{"status": STATUS_DELETE}})
	c.JSON(200, gin.H{"code": 200, "spell": "DELETE"})
}
