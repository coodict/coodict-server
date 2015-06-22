/*
PACKAGE routes
rainy @ 2015-06-08 <me@rainy.im>
*/
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

func (app *App) githubAuth(c *gin.Context) {
	code := c.Query("code")
	m := map[string]interface{}{
		"code":          code,
		"client_id":     GITHUB_CLIENT,
		"client_secret": GITHUB_SECRET,
	}
	mJson, _ := json.Marshal(m)

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(mJson))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	var dat map[string]string
	if err := json.Unmarshal(body, &dat); err != nil {
		panic(err)
	}
	fmt.Println(dat["access_token"])
	req, err = http.NewRequest("GET", "https://api.github.com/user", nil)
	req.Header.Set("Authorization", "token "+dat["access_token"])
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err = client.Do(req)
	defer resp.Body.Close()

	var githubUser map[string]interface{}
	body, err = ioutil.ReadAll(resp.Body)
	if err = json.Unmarshal(body, &githubUser); err != nil {
		panic(err)
	}
	fmt.Printf("%+v", githubUser)
	var u User
	app.db.C("user").Find(bson.M{"mail": githubUser["email"].(string)}).One(&u)
	if u.Name == "" {
		user := User{
			Name:  githubUser["login"].(string),
			Mail:  githubUser["email"].(string),
			Salt:  "DEFAULT_SALT",
			Pass:  "GITHUB_WITHOUT_PASS",
			Coins: 100,
			Votes: 0,
		}
		err = app.db.C("user").Insert(user)
		if err != nil {
			//c.JSON(200, gin.H{"code": 500, "msg": "服务器错误，请稍后再试！"})
			c.Redirect(301, "http://localhost:9000/#/p500/服务器错误，请稍后再试！")
			return
		}
		u = user
	}
	tokenString, _ := genUsrToken(u)
	redirUrl := fmt.Sprintf("http://localhost:9000/#/callback/%s", tokenString)
	c.Redirect(301, redirUrl)
}

func (app *App) signup(c *gin.Context) {
	var signup Signup
	user := new(User)
	err := c.Bind(&signup)
	if err != nil {
		// Check required bindings
		c.JSON(200, gin.H{"code": 402, "msg": "参数不完整！"})
	} else {
		// Check promo code
		var promo Promo
		app.db.C("promo").Find(bson.M{"code": signup.Promo, "isused": 0}).One(&promo)
		if promo.Code == "" {
			c.JSON(200, gin.H{"code": 400, "msg": "无效的邀请码！"})
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
				c.JSON(200, gin.H{"code": 500, "msg": "服务器错误，请稍后再试！"})
				return
			}
			tokenString, _ := genUsrToken(u)
			c.JSON(200, gin.H{"code": 200, "msg": "OK", "jwt": tokenString})
		} else {
			c.JSON(200, gin.H{"code": 401, "msg": "邮箱/ID已存在！"})
		}
	}
}

func (app *App) signin(c *gin.Context) {
	var signin Signin
	err := c.Bind(&signin)
	if err != nil {
		c.JSON(400, gin.H{"code": 402, "msg": "Not enough params!"})
		return
	}
	var user User
	app.db.C("user").Find(bson.M{"$or": []bson.M{
		{"name": signin.Name},
		{"mail": signin.Name},
	}}).One(&user)
	if user.Name == "" {
		c.JSON(200, gin.H{"code": 400, "msg": "用户名不存在！"})
		return
	}
	if genMd5(user.Salt+signin.Pass) != user.Pass {
		c.JSON(200, gin.H{"code": 401, "msg": "密码错误！"})
		return
	}
	tokenString, _ := genUsrToken(user)
	c.JSON(200, gin.H{"code": 200, "msg": "OK", "jwt": tokenString, "tgs": user.Tags, "lgs": user.Lang})
}
