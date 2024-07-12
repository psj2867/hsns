package models_test

import (
	"testing"

	"github.com/psj2867/hsns/models"
	"github.com/stretchr/testify/assert"
)

func TestAddUser(t *testing.T) {
	user := models.User{
		UserId: "1234",
		Name:   "asdf",
	}
	err := user.Add()
	if err != nil {
		t.Error(err)
		return
	}

	user2 := models.User{}
	user2.Get(user.Id)
	assert.Equal(t, user.Id, user2.Id)
}
