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

// rpc
type RegisterArgs struct{}

type RegisterReply struct {
	Id int
}

type DispatchArgs struct {
	Id int
}

type DispatchReply struct {
	Task int
	Flag string
	File string
}

type RunArgs struct {
	Id   int
	Flag string
	File string
}

type RunReply struct{}

type FinishArgs struct {
	Id   int
	Task int
}

type FinishReply struct{}

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the coordinator.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func coordinatorSock() string {
	s := "/var/tmp/5840-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
