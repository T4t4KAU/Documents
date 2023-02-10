package mr

import (
	"fmt"
	"io"
	"log"
	"net/rpc"
	"os"
	"strconv"
)

// Map functions return a slice of KeyValue.
type KeyValue struct {
	Key   string
	Value string
}

func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {
	// get worker id
	var worker int = register()
	for {
		// ask coordinator for tasks
		id, t, ok := ask(worker)
		if t.file == "" {
			break
		}
		if ok {
			// accept task and start
			start(t, worker)
			text := process(t.file)
			switch t.flag {
			case "map":
				kva := mapf(t.file, text)
				save(kva, id, worker) // save result in file
			case "reduce":
			}
			// notify coordinator
			finish(id, t, worker)
		}
	}
}

func register() int {
	reply := RegisterReply{}
	ok := call("Coordinator.Register", &RegisterArgs{}, &reply)
	if ok {
		return reply.Id
	} else {
		fmt.Printf("register: call failed!\n")
	}
	return 0
}

func ask(id int) (int, task, bool) {
	args := DispatchArgs{Id: id}
	reply := DispatchReply{}
	ok := call("Coordinator.Dispatch", &args, &reply)
	if ok {
		if reply.File == "" {
			return 0, task{}, false
		}
		fmt.Printf("worker %d get %s task: %v\n", id, reply.Flag, reply.File)
		return reply.Task, task{reply.Flag, reply.File}, true
	} else {
		fmt.Printf("ask: call failed!\n")
		return 0, task{}, false
	}
}

func start(t task, id int) {
	args := RunArgs{Flag: t.flag, File: t.file, Id: id}
	ok := call("Coordinator.Run", &args, &RunReply{})
	if ok {
		fmt.Printf("%s task: %s start\n", args.Flag, args.File)
	} else {
		fmt.Printf("start: call failed!\n")
	}
}

func finish(id int, t task, worker int) {
	args := FinishArgs{worker, id}
	reply := RunReply{}
	ok := call("Coordinator.Finish", &args, &reply)
	if ok {
		fmt.Printf("%s task %s finished\n", t.flag, t.file)
	} else {
		fmt.Printf("finish: call failed!\n")
	}
}

func call(rpcname string, args interface{}, reply interface{}) bool {
	sockname := coordinatorSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}

func process(file string) string {
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("cannot open %s\n", file)
	}
	content, err := io.ReadAll(f)
	if err != nil {
		log.Fatalf("cannot read %s\n", file)
	}
	f.Close()
	return string(content)
}

func save(kva []KeyValue, id int, worker int) {
	oname := "mr-out-" + strconv.Itoa(id)
	f, err := os.OpenFile(oname, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("worker %d open file: %s error: %v\n", id, oname, err)
	}
	defer f.Close()

	for _, kv := range kva {
		fmt.Fprintf(f, "%s %s\n", kv.Key, kv.Value)
	}
	fmt.Printf("worker %d has saved file: %s\n", worker, oname)
}
