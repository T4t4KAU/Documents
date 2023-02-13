package mr

//
// RPC definitions.
//
// remember to capitalize all names.
//

import (
	"os"
	"strconv"
)

//
// example to show how to declare the arguments
// and reply for an RPC.
//

type ExampleArgs struct {
	X int
}

type ExampleReply struct {
	Y int
}

// Add your RPC definitions here.
type RegistrArgs struct{}

type RegisterReply struct {
	Id      int32
	NReduce int
}

type AssignArgs struct {
	Id int32 // worker id
}

type AssignReply struct {
	Id   int32 // task id
	File string
	Flag string
}

type RoundArgs struct {
	WID  int32
	TID  int32
	Flag string
}

type RoundReply struct {
}

type FinishArgs struct {
	WID  int32
	TID  int32
	Flag string
}

type FinishReply struct {
	Accept bool
}

type ReviewArgs struct {
	WID int32
	TID int32
}

type ReviewReply struct{}

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the coordinator.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func coordinatorSock() string {
	s := "/var/tmp/5840-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
