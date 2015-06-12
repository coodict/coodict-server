/*
rainy @ 2015-06-05 <me@rainy.im>
*/
package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

var MOGODB_URI = "mongodb://127.0.0.1:27017/cobble"

type App struct {
	g  *gin.Engine
	db *mgo.Database
}

func main() {
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		// Run this on all requests
		// Should be moved to a proper middleware
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Next()
	})
	router.OPTIONS("/*cors", func(c *gin.Context) {
		// Empty 200 response
	})
	sess, err := mgo.Dial(MOGODB_URI)
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}
	defer sess.Close()
	app := &App{g: router, db: sess.DB("cobble")}

	router.Static("/brace/mode", "./brace/mode/")
	userAPI := router.Group("/user")
	{
		userAPI.POST("/signup", app.signup)
		userAPI.POST("/signin", app.signin)
	}
	spellAPI := router.Group("/spell")
	spellAPI.Use(Auth(mySigningKey))
	{
		spellAPI.POST("/create", app.createSpell)

	}
	router.POST("/fetchSpell", app.fetchSpell)
	router.Run(":8080")
}
