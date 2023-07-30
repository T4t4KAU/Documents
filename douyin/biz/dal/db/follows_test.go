package db

import (
	"douyin/biz/mw/redis"
	"fmt"
	"testing"
)

func TestAddNewFollow(t *testing.T) {
	Init()
	redis.Init()

	f := Follows{
		UserId:     1001,
		FollowerId: 1004,
	}

	ok, err := AddNewFollow(&f)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	if ok {
		fmt.Println("success")
	} else {
		fmt.Println("failed")
	}
}

func TestDeleteFollow(t *testing.T) {
	Init()
	redis.Init()

	f := Follows{
		UserId:     1009,
		FollowerId: 1006,
	}

	ok, err := DeleteFollow(&f)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	if ok {
		fmt.Println("success")
	} else {
		fmt.Println("failed")
	}
}

func TestGetFollowCount(t *testing.T) {
	Init()
	redis.Init()

	cnt, err := GetFollowerCount(1009)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	fmt.Println(cnt)
}

func TestGetFollowerList(t *testing.T) {
	Init()
	redis.Init()

	followerList, err := getFollowIdList(1001)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	for _, fid := range followerList {
		fmt.Println(fid)
	}
}
