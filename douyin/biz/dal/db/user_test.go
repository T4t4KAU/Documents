package db

import (
	"fmt"
	"testing"
)

func TestCreateUser(t *testing.T) {
	Init()
	user := &User{
		ID:       1005,
		UserName: "test",
		Password: "123456",
	}

	uid, err := CreateUser(user)
	if err != nil {
		t.Errorf(err.Error())
	}

	fmt.Printf("%v\n", uid)
}

func TestQueryUserByName(t *testing.T) {
	Init()
	user, err := QueryUserByName("test")
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	fmt.Printf("%v\n", user)
}

func TestQueryUserById(t *testing.T) {
	Init()
	user, err := QueryUserById(int64(1))
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	fmt.Println(user)
}
