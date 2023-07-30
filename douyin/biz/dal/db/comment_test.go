package db

import (
	"fmt"
	"testing"
)

func TestAddNewComment(t *testing.T) {
	Init()
	comment := &Comment{
		UserId:      1000,
		VideoId:     115,
		CommentText: "video comment test",
	}

	err := AddNewComment(comment)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	fmt.Println("insert successfully")
}

func TestDeleteCommentById(t *testing.T) {
	Init()
	err := DeleteCommentById(int64(2))
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	fmt.Println("delete successfully")
}
