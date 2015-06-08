/*
PACKAGE helpers
rainy @ 2015-06-08 <me@rainy.im>
*/
package main

import (
	"crypto/md5"
	"encoding/hex"
)

// Helper funcs
func genSalt(name string) string {
	return genMd5(name)[:8]
}
func genMd5(txt string) string {
	hasher := md5.New()
	hasher.Write([]byte(txt))
	return hex.EncodeToString(hasher.Sum(nil))
}
