package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPersistence(t *testing.T) {
	// <setup code>

	// var err = os.Remove("unfollow_test.db")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	fmt.Println("==> done deleting file")

	t.Run("InsertUser", func(t *testing.T) {
		id := int64(123)
		us := User{ID: id}
		us.Save()

		ug := &User{ID: id}
		ug.Get()
		fmt.Println(ug)

		assert.Equal(t, id, ug.ID)

	})
	t.Run("UpdateUser", func(t *testing.T) {
		ug := &User{ID: int64(123)}
		ug.Get()
		screenName := "philipsahli"
		ug.ScreenName = screenName
		ug.Save()

		ur := &User{ID: int64(123)}
		ur.Get()
		assert.Equal(t, screenName, ur.ScreenName)

	})
	t.Run("TestDeleteUser", func(t *testing.T) {
		ug := &User{ID: int64(123)}
		ug.Get()
		ug.Delete()

	})

	var err = os.Remove("unfollow_test.db")
	if err != nil {
		fmt.Println(err.Error())
	}

}
