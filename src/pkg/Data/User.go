package Data

import (
	"launchpad.net/gobson/bson"
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
	if i != 0 {
		return -1, TakenUser
	}
	
	if e != nil {
		panic(e)
	}

	i, e = CUsers.Find(bson.M{"email": strings.ToLower(email)}).Count()
	if i != 0 {
		return -2, TakenEmail
	} 
	
	if e != nil {
		panic(e)
	}

	return 1, OK
}

func RegisterUser(user *User) bool {
	if i, _ := CheckUser(user.User, user.EMail, user.Password); i == 1 {
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
		panic(e)
	}
	if u.User == "" {
		return -3, BadUser, ""
	}

	if u.Password != password {
		return -4, BadPassword, ""
	}

	return 0, LoginOK, u.ID
}
