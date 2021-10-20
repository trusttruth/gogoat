package utils

import (
	"bytes"
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
)

func GetSha256(str string) string {
	hash := crypto.SHA256.New()
	hash.Write([]byte(str))
	hashpass := hash.Sum(nil)
	return hex.EncodeToString(hashpass)
}

func GetRandMd5Sum() string {
	crutime := time.Now().Unix()
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(crutime, 10))
	token := fmt.Sprintf("%x", h.Sum(nil))
	// token := hex.EncodeToString(h.Sum(nil))
	return token
}

func GetRedisConn() redis.Conn {
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("redis dial err")
		return nil
	}
	return conn
}

func GetRandString(len int) string {
	var container string
	var str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0; i < len; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		container += string(str[randomInt.Int64()])
	}
	return container
}
func checkPassReg(pass string) bool {
	if len(pass) < 6 {
		fmt.Errorf("BAD PASSWORD: The password is shorter than 6 characters")
		return false
	}
	if len(pass) > 12 {
		fmt.Errorf("BAD PASSWORD: The password is logner than 12 characters")
		return false
	}
	level := 0
	patternList := []string{`[0-9]+`, `[a-z]+`, `[A-Z]+`, `[~!@#$%^&*?_-]+`}
	for _, pattern := range patternList {
		match, _ := regexp.MatchString(pattern, pass)
		if match {
			level++
		}
	}

	if level < 3 {
		fmt.Errorf("The password does not satisfy the current policy requirements. ")
		return false
	}
	return true
}

func GetUsername(r *http.Request) string {
	sid, err := r.Cookie("sessionid")
	if err != nil {
		return ""
	}
	if sid.Value != "" {
		redisconn := GetRedisConn()
		if redisconn == nil {
			return ""
		}
		defer redisconn.Close()
		v, err := redisconn.Do("get", sid.Value)

		if err != nil {
			fmt.Println(err.Error())
		}
		// fmt.Println(v.(string))
		// return v.(string)
		// return "test"
		return string(v.([]uint8))
	}
	return ""

}

func Validate(inter interface{}, tp string) bool {
	switch inter.(type) {
	case string:
		if tp == "username" {
			return len(inter.(string)) > 6
		} else {
			if tp == "password" {
				return checkPassReg(inter.(string))

			}
		}
	default:
		return false
	}
	return false
}
