package db

import (
	"testing"
	"time"
)

func TestAddNewMessage(t *testing.T) {
	Init()
	message := &Messages{
		ToUserId:   1001,
		FromUserId: 1004,
		Content:    "test: 1004 -> 1001, this is message",
	}

	ok, err := AddNewMessage(message)
	if err != nil {
		t.Errorf("Add message 1 error: %v", err)
	}
	if !ok {
		t.Logf("Failed to add message 1")
	}

	time.Sleep(time.Second)
	message = &Messages{
		ToUserId:   1004,
		FromUserId: 1001,
		Content:    "test: 1001 -> 1004, this is latest message",
	}

	ok, err = AddNewMessage(message)
	if err != nil {
		t.Fatalf("Add message 2 error: %v", err)
	}
	if !ok {
		t.Logf("Failed to add message 2")
	}
}
