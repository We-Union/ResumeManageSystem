package utils

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func GetToken(length int) (token string) {
	rand.Seed(time.Now().UnixNano())
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	var result []byte
	for i := 0; i < length; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}

	return string(result)
}

func GetUidInt(c *gin.Context) (uidInt int) {
	session := sessions.Default(c)
	uid := session.Get("uid")

	if uid == nil {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您还未登录，请先登录"})
		return -1
	}
	uidInt_ := uid.(int)
	return uidInt_
}
func IsExist(path string) (exist bool) {
	_, err := os.Stat(path)
	return err == nil
}
