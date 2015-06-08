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
