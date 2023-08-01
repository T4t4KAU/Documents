package rpc

import (
	"context"
	"douyin/kitex_gen/relation"
	"fmt"
	"testing"
)

func TestRelationAction(t *testing.T) {
	InitRPC()

	resp, err := RelationAction(context.Background(), &relation.RelationActionRequest{
		CurrentUserId: 1001,
		ToUserId:      1002,
		ActionType:    1,
	})

	if err != nil {
		t.Errorf(err.Error())
	}

	fmt.Printf("%#v\n", resp)
}
