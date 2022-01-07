package twitter

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm/clause"

	"gorm.io/gorm"
)

var err error
var db *gorm.DB

func init() {
	//dbName := "unfollow.db"
	dbName := "unfollow.db?cache=shared&mode=rwc&busy_timeout=200"
	if flag.Lookup("test.v") != nil || dbName == "" {
		//dbName = "unfollow.db?cache=shared&mode=rwc&busy_timeout=200"
		//dbName = "unfollow.db"
		dbName = "../../unfollow_test.db?cache=shared&mode=rwc&busy_timeout=200"
	}
	log.Println("opening", dbName)
	log.Println(os.Getwd())
	db, err = gorm.Open(sqlite.Open(dbName), &gorm.Config{})

	if err != nil {
		log.Fatal()
	}

	err = db.AutoMigrate(&User{})
	if err != nil {
		fmt.Println(err)
		log.Fatal()
	}
	err = db.AutoMigrate(&Follower{})
	if err != nil {
		fmt.Println(err)
		log.Fatal()
	}
	println("init db", db)
	// defer db.Close()
}

func (u *User) Save() error {
	if err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},                     // key colume
		DoUpdates: clause.AssignmentColumns([]string{"screen_name"}), // column needed to be updated
	}).Create(&u).Error; err != nil {

		return err
	}

	return nil
}

func (u *User) Get() error {
	err := db.First(&u, u.ID).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Delete() {
	db.Delete(&u)
}

func GetUsers() []User {
	users := []User{}
	tx := db.Find(&users)
	if tx.Error != nil {
		log.Fatal(tx)
	}

	return users
}
