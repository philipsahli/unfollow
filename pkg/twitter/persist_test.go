package twitter

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPersistence(t *testing.T) {
	// <setup code>

	var err = os.Remove("../../unfollow_test.db")
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("==> done deleting file")

	t.Run("InsertUser", func(t *testing.T) {
		id := int64(123)
		us := User{ID: id}
		err := us.Save()
		assert.Nil(t, err)

		ug := &User{ID: id}
		err = ug.Get()
		assert.Nil(t, err)

		assert.Equal(t, id, ug.ID)

	})

	t.Run("UpdateUser", func(t *testing.T) {
		ug := &User{ID: int64(123)}
		err := ug.Get()
		assert.Nil(t, err)
		ug.ScreenName = "philipsahli"
		err = ug.Save()
		assert.Nil(t, err)

		ur := &User{ID: int64(123)}
		err = ur.Get()
		assert.Nil(t, err)
		assert.Equal(t, ug.ScreenName, ur.ScreenName)

	})

	t.Run("TestDeleteUser", func(t *testing.T) {
		ug := &User{ID: int64(123)}
		ug.Get()
		ug.Delete()

	})

	err = os.Remove("unfollow_test.db")
	if err != nil {
		fmt.Println(err.Error())
	}

}
