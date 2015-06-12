/*
PACKAGE routes
rainy @ 2015-06-08 <me@rainy.im>
*/
package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

func (app *App) signup(c *gin.Context) {
	var signup Signup
	user := new(User)
	err := c.Bind(&signup)
	if !err {
		// Check required bindings
		c.JSON(200, gin.H{"code": 402, "msg": "Params are Required"})
	} else {
		// Check promo code
		var promo Promo
		app.db.C("promo").Find(bson.M{"code": signup.Promo, "isused": 0}).One(&promo)
		if promo.Code == "" {
			c.JSON(200, gin.H{"code": 400, "msg": "Promocode not available"})
			return
		}
		app.db.C("user").Find(bson.M{"$or": []bson.M{{"name": signup.Name}, {"mail": signup.Mail}}}).One(&user)
		if user.Name == "" {
			// Create a new user
			salt := genSalt(signup.Name)
			d := time.Now().Local()

			u := User{
				Name:       signup.Name,
				Mail:       signup.Mail,
				Salt:       salt,
				Pass:       genMd5(salt + signup.Pass),
				Coins:      100,
				Spells:     0,
				Votes:      0,
				CreateDate: fmt.Sprintf("%d-%d-%d", d.Year(), d.Month(), d.Day()),
			}
			err := app.db.C("user").Insert(u)
			if err != nil {
				c.JSON(200, gin.H{"code": 500, "msg": "Our fault"})
				return
			}
			tokenString, _ := genUsrToken(u)
			c.JSON(200, gin.H{"code": 200, "msg": "OK", "jwt": tokenString})
		} else {
			c.JSON(200, gin.H{"code": 401, "msg": "Already Exist"})
		}
	}
}

func (app *App) signin(c *gin.Context) {
	var signin Signin
	err := c.Bind(&signin)
	if !err {
		c.JSON(400, gin.H{"code": 402, "msg": "Not enough params!"})
		return
	}
	var user User
	app.db.C("user").Find(bson.M{"$or": []bson.M{
		{"name": signin.Name},
		{"mail": signin.Name},
	}}).One(&user)
	fmt.Println(user)
	if user.Name == "" {
		c.JSON(400, gin.H{"code": 400, "msg": "No such user"})
		return
	}
	if genMd5(user.Salt+signin.Pass) != user.Pass {
		c.JSON(200, gin.H{"code": 401, "msg": "wrong password"})
		return
	}
	tokenString, _ := genUsrToken(user)
	c.JSON(200, gin.H{"code": 200, "msg": "OK", "jwt": tokenString})
}