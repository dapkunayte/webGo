package models

import (
	"errors"
	//"time"
)

var ErrNoRecord = errors.New("models: подходящей записи не найдено")

var ErrNoUser = errors.New("users: нет такого пользователя")

var ErrNoSub = errors.New("subsribes: нет такой подписки")

type Note struct {
	ID       int
	Title    string
	Content  string
	Created  string
	Username string
}

type User struct {
	Username string `json:"username", db:"username"`
	Password string `json:"password", db:"password"`
	Email    string `json:"email", db:"email"`
}

type Subscribe struct {
	ID       int
	Date     string
	SubId    string
	FollowId string
}
