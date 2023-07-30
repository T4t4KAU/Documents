package db

import (
	"fmt"
	"testing"
	"time"
)

func TestGetVideosByLastTime(t *testing.T) {
	Init()
	lastTime := time.Now()
	videos, err := GetVideosByLastTime(lastTime)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	for _, video := range videos {
		fmt.Printf("%#v\n", video)
	}
}

func TestGetVideosByUserID(t *testing.T) {
	Init()
	uid := int64(1000)
	videos, err := GetVideoByUserId(uid)
	if err != nil {
		t.Errorf(err.Error())
	}
	for _, video := range videos {
		fmt.Printf("%#v\n", video)
	}
}
