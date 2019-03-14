package main

import (
	"flag"
	"fmt"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB
var err error

func init() {
	var dbName string
	if flag.Lookup("test.v") == nil {
		dbName = "unfollow.db"
	} else {
		dbName = "unfollow_test.db"
	}
	db, err = gorm.Open("sqlite3", dbName)

	db.CreateTable(&User{})
	db.CreateTable(&Follower{})
	println("init db", db)
	if err != nil {
		fmt.Println(err)
	}
	// defer db.Close()
}

func (u *User) Save() {
	if db.NewRecord(u) {
		db.Create(&u)
		fmt.Println("User created", u.ID, u.ScreenName)
	} else {
		db.Save(&u)
		fmt.Println("User updated", u.ID, u.ScreenName)
	}
}

func (u *User) Get() {
	db.First(&u, u.ID)
}

func (u *User) Delete() {
	db.Delete(&u)
}
