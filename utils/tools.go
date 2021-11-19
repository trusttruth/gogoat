package utils

import (
	"bytes"
	"context"
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"regexp"
	"strconv"
	"time"

	// "github.com/gomodule/redigo/redis"
	"database/sql"

	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/configor"

	_ "github.com/go-sql-driver/mysql"
)

var ctx = context.Background()

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

// func GetRedisConn() redis.Conn {
// 	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
// 	if err != nil {
// 		fmt.Println("redis dial err")
// 		return nil
// 	}
// 	return conn
// }
func GetRedisConn() *redis.Client {
	conn := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
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
			fmt.Println(err.Error())
			return ""
		}
		defer redisconn.Close()
		v, err := redisconn.Get(ctx, sid.Value).Result()

		if err != nil {
			fmt.Println(err.Error())
		}
		return v
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

func GetDBconnStr() string {
	configor.Load(&Config, "pkg/vul/config.yaml")
	dbcon := Config.DB.User + ":" + Config.DB.Password + "@tcp(" + Config.DB.Host + ":" + Config.DB.Port + ")/" + Config.DB.Name + "?charset=utf8"
	return dbcon
}

func GetUserpassFromDb(username string) string {
	var pass string
	_, err := X.Table("user").Where("username = ?", username).Cols("password").Get(&pass)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return pass
}

func ChangePass(pass, username string) {
	user := &User{Username: username}
	ok, _ := X.Get(user)
	if ok {
		X.ID(user.Id).Omit("username").Update(&User{Password: pass})
	}
}

func GetRawDBcon() *sql.DB {
	db, err := sql.Open("mysql", GetDBconnStr())
	if err != nil {
		fmt.Println("raw connect db err")
		panic("raw connect db err")
	}
	return db
}

func Clearcomment(w http.ResponseWriter, r *http.Request) {
	sql := "truncate  TABLE  message;"
	_, err := X.Exec(sql)
	if err != nil {
		fmt.Println(err)
	}
	return

}

//tran struct src to struct dst
func TranStruct(src interface{}, dst interface{}) error {
	s1, err := json.Marshal(src)
	if err != nil {
		return err
	}
	// fmt.Println(string(s1))
	err = json.Unmarshal(s1, dst)
	// fmt.Println(dst)
	if err != nil {
		return err
	}
	return nil
}
