/*
rainy @ 2015-06-05 <me@rainy.im>
*/
package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var MOGODB_URI = "mongodb://127.0.0.1:27017/cobble"

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

type App struct {
	g  *gin.Engine
	db *mgo.Database
}

func main() {
	router := gin.Default()
	sess, err := mgo.Dial(MOGODB_URI)
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}
	defer sess.Close()
	app := &App{g: router, db: sess.DB("cobble")}

	userAPI := router.Group("/user")
	{
		userAPI.POST("/signup", app.signup)
		userAPI.POST("/signin", app.signin)
	}

	router.Run(":8080")
}

func (app *App) signup(c *gin.Context) {
	var signup Signup
	user := new(User)
	err := c.Bind(&signup)
	if !err {
		// Check required bindings
		c.JSON(400, gin.H{"code": 402, "msg": "Params are Required"})
	} else {
		// Check promo code
		var promo Promo
		app.db.C("promo").Find(bson.M{"code": signup.Promo, "isused": 0}).One(&promo)
		if promo.Code == "" {
			c.JSON(400, gin.H{"code": 400, "msg": "Promocode not available"})
			return
		}
		app.db.C("user").Find(bson.M{"$or": []bson.M{{"name": signup.Name}, {"mail": signup.Mail}}}).One(&user)
		if user.Name == "" {
			// Create a new user
			salt := genSalt(signup.Name)
			u := User{
				Name:   signup.Name,
				Mail:   signup.Mail,
				Salt:   salt,
				Pass:   genMd5(salt + signup.Pass),
				Coins:  100,
				Spells: 0,
				Votes:  0,
			}
			err := app.db.C("user").Insert(u)
			if err != nil {
				c.JSON(500, gin.H{"code": 500, "msg": "Our fault"})
				return
			}
			c.JSON(200, gin.H{"code": 200, "msg": "OK", "user": gin.H{
				"name":   u.Name,
				"mail":   u.Mail,
				"coins":  u.Coins,
				"spells": u.Spells,
				"votes":  u.Votes,
			}})
		} else {
			c.JSON(400, gin.H{"code": 401, "msg": "Already Exist"})
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
	fmt.Println(signin)
	if user.Name == "" {
		c.JSON(400, gin.H{"code": 400, "msg": "No such user"})
		return
	}
	if genMd5(user.Salt+signin.Pass) != user.Pass {
		c.JSON(200, gin.H{"code": 401, "msg": "wrong password"})
		return
	}
	c.JSON(200, gin.H{"code": 200, "msg": "OK", "user": gin.H{
		"name":   user.Name,
		"mail":   user.Mail,
		"coins":  user.Coins,
		"spells": user.Spells,
		"votes":  user.Votes,
	}})
}

// Helper funcs
func genSalt(name string) string {
	return genMd5(name)[:8]
}
func genMd5(txt string) string {
	hasher := md5.New()
	hasher.Write([]byte(txt))
	return hex.EncodeToString(hasher.Sum(nil))
}
