package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/michaelvanstraten/prometheus/database"
)

type User struct {
	Username string `redis:"username"`
	Password string `redis:"password"`
	Test string `redis:"fett"`
}

func main() {
	
	var db = redis.NewClient(
		&redis.Options{
			Addr: "127.0.0.1:6379",
		},
	)
	var users = database.NewCollection("user", db)

	var user = make([]User, 3)
	var p = make([]interface{}, 3)
	for i := range user {
		p[i] = &user[i]
	}
	users.GetMembers(p, users.Sets["MASTER"], 0, "username")
	for u := range user {
		println(user[u].Username)
	}
}