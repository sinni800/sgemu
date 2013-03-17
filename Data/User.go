package Data

import (
	"crypto/md5"
	"encoding/hex"
	"labix.org/v2/mgo/bson"
	"strings"
)

type User struct {
	ID       string "_id"
	User     string
	EMail    string
	Password string
}

const (
	TakenUser   = "This username is already taken."
	TakenEmail  = "This email is already in use."
	OK          = "New user has been created!"
	BadUser     = "This user is not exist!"
	BadPassword = "Your password is invalid"
	LoginOK     = "Good, now wait for the next rev :]"
)

func CheckUser(user, email, password string) (int, string) {
	i, e := CUsers.Find(bson.M{"user": strings.ToLower(user)}).Count()
	if e != nil {
		panic(e)
	}

	if i != 0 {
		return -1, TakenUser
	}

	i, e = CUsers.Find(bson.M{"email": strings.ToLower(email)}).Count()
	if e != nil {
		panic(e)
	}
	if i != 0 {
		return -2, TakenEmail
	}

	return 1, OK
}

func RegisterUser(user *User) bool {
	if i, _ := CheckUser(user.User, user.EMail, user.Password); i == 1 {
		MD5 := md5.New()
		MD5.Write([]byte(user.Password))

		user.Password = hex.EncodeToString(MD5.Sum(nil))

		user.User = strings.ToLower(user.User)
		user.EMail = strings.ToLower(user.EMail)
		CUsers.Insert(user)
		return true
	}
	return false
}

func Login(user, password string) (int, string, string) {
	u := User{}
	e := CUsers.Find(bson.M{"user": strings.ToLower(user)}).One(&u)
	if e != nil {
		return -3, BadUser, ""
	}
	if u.User == "" {
		return -3, BadUser, ""
	}

	MD5 := md5.New()
	MD5.Write([]byte(password))

	if u.Password != hex.EncodeToString(MD5.Sum(nil)) {
		return -4, BadPassword, ""
	}

	return 0, LoginOK, u.ID
}
