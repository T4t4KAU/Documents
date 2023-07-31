package db

import (
	"context"
	"fmt"
	"testing"
)

func TestCreateUser(t *testing.T) {
	Init()
	user := &User{
		ID:       1000,
		UserName: "test",
		Password: "123456",
	}

	uid, err := CreateUser(context.Background(), user)
	if err != nil {
		t.Errorf(err.Error())
	}

	fmt.Printf("%v\n", uid)
}

func TestQueryUserByName(t *testing.T) {
	Init()
	user, err := QueryUserByName(context.Background(), "test")
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	fmt.Printf("%v\n", user)
}

func TestQueryUserById(t *testing.T) {
	Init()
	user, err := QueryUserById(context.Background(), int64(1000))
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	fmt.Println(user)
}
