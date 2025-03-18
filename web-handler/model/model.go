package model

import (
	"fmt"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	LastName  string    `json:"last_name"`
	FirstName string    `json:"first_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func NewUser() *User {
	return &User{}
}

var users map[int]User = make(map[int]User)
var check map[string]int = map[string]int{} // user.email -> user.id
var index = 0

func CreateUser(newUser User) error {
	if _, exist := check[newUser.Email]; !exist {
		index++
		newUser.ID = index
		users[index] = newUser
		check[newUser.Email] = index
		return nil
	}
	return fmt.Errorf("already exist user, user_email=%s", newUser.Email)
}

func GetUserById(id int) (User, error) {
	if user, exist := users[id]; exist {
		if user.ID == 0 {
			return User{}, fmt.Errorf("not exist %d user", id)
		}
		return user, nil
	}
	return User{}, fmt.Errorf("not exist %d user", id)
}

func DeleteUserByEmail(delUser User) error {
	if id, exist := check[delUser.Email]; exist {
		delete(users, id)
		delete(check, delUser.Email)
		return nil
	}
	return fmt.Errorf("not exist %s user", delUser.Email)
}

func UpdateUser(updateUser User) error {
	if id, exist := check[updateUser.Email]; exist {
		users[id] = updateUser
		return nil
	}
	return fmt.Errorf("not exist %s user", updateUser.Email)
}

func PutUser(putUser User) (User, error) {
	if id, exist := check[putUser.Email]; exist {
		user := users[id]
		if putUser.FirstName != "" {
			user.FirstName = putUser.FirstName
		}
		if putUser.LastName != "" {
			user.LastName = putUser.LastName
		}
		users[id] = user
		return user, nil
	}

	return User{}, fmt.Errorf("not exist %s user", putUser.Email)
}
