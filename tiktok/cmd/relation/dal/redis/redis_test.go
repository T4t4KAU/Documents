package redis

import (
	"context"
	"fmt"
	"testing"
)

func TestAddFollow(t *testing.T) {
	Init()

	ctx := context.Background()
	f := Follows{}
	f.AddFollow(ctx, 1, 2)
}

func TestCheckFollow(t *testing.T) {
	Init()

	ctx := context.Background()
	f := Follows{}
	ok := f.CheckFollow(ctx, 2)
	fmt.Println(ok)
}

func TestDelFollow(t *testing.T) {
	Init()

	ctx := context.Background()
	f := Follows{}
	f.DelFollow(ctx, 1, 2)

	ok := f.CheckFollow(ctx, 2)
	if ok {
		fmt.Println("failed to delete")
	}
	fmt.Println("success")
}

func TestExistFollow(t *testing.T) {
	Init()

	ctx := context.Background()
	f := Follows{}
	ok := f.ExistFollow(ctx, 1, 2)
	fmt.Println(ok)
}

func TestCountFollow(t *testing.T) {
	Init()

	ctx := context.Background()
	f := Follows{}
	cnt, err := f.CountFollow(ctx, 2)
	if err != nil {
		return
	}
	fmt.Println(cnt)
}

func TestAddFollower(t *testing.T) {
	Init()

	ctx := context.Background()
	f := Follows{}
	f.AddFollower(ctx, 1, 2)
}

func TestCheckFollower(t *testing.T) {
	Init()

	ctx := context.Background()
	f := Follows{}
	ok := f.CheckFollower(ctx, 1)
	fmt.Println(ok)
}

func TestDelFollower(t *testing.T) {
	Init()

	ctx := context.Background()
	f := Follows{}
	f.DelFollower(ctx, 1, 2)

	ok := f.CheckFollow(ctx, 1)
	if ok {
		fmt.Println("failed to delete")
	}
	fmt.Println("success")
}

func TestExistFollower(t *testing.T) {
	Init()

	ctx := context.Background()
	f := Follows{}
	ok := f.ExistFollower(ctx, 1, 2)
	fmt.Println(ok)
}
