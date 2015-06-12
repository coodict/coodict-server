/*
PACKAGE helpers
rainy @ 2015-06-08 <me@rainy.im>
*/
package main

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var mySigningKey = "AVATQ!#@$#^%ASBA1354"

// Helper funcs
func genSalt(name string) string {
	return genMd5(name)[:8]
}
func genMd5(txt string) string {
	hasher := md5.New()
	hasher.Write([]byte(txt))
	return hex.EncodeToString(hasher.Sum(nil))
}
func genUsrToken(u User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	// Headers
	token.Header["alg"] = "HS256"
	token.Header["typ"] = "JWT"

	// Claims
	token.Claims["name"] = u.Name
	token.Claims["mail"] = u.Mail
	token.Claims["coins"] = u.Coins
	token.Claims["spells"] = u.Spells
	token.Claims["votes"] = u.Votes
	token.Claims["date"] = u.CreateDate

	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	tokenString, err := token.SignedString([]byte(mySigningKey))

	return tokenString, err
}
