package mr

import (
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"net/rpc"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Map functions return a slice of KeyValue.
type KeyValue struct {
	Key   string
	Value string
}

type ByKey []KeyValue

func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }

// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

var (
	wid int32 // worker id
	nr  int   // nreduce
)

type mf func(string, string) []KeyValue
type rf func(string, []string) string

// main/mrworker.go calls this function.
func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {

	wid, nr = register()
	if wid <= 0 {
		log.Printf("invalid worker id\n")
	}
	log.Printf("worker id: %d\n", wid)
	for {
		tid, task := ask(wid)
		if tid == -SUSPEND || tid == -READY {
			log.Printf("wait for task...\n")
			time.Sleep(10 * time.Millisecond)
			continue
		}
		if tid == -DONE {
			log.Printf("tasks all done\n")
			return
		}
		log.Printf("get %s task: %s id: %d\n", task.Flag, task.File, tid)
		start(wid, tid, task.Flag)
		work(tid, task, mapf, reducef)
		done(tid)
	}
}

// example function to show how to make an RPC call to the coordinator.
//
// the RPC argument and reply types are defined in rpc.go.
func CallExample() {

	// declare an argument structure.
	args := ExampleArgs{}

	// fill in the argument(s).
	args.X = 99

	// declare a reply structure.
	reply := ExampleReply{}

	// send the RPC request, wait for the reply.
	// the "Coordinator.Example" tells the
	// receiving server that we'd like to call
	// the Example() method of struct Coordinator.
	ok := call("Coordinator.Example", &args, &reply)
	if ok {
		// reply.Y should be 100.
		fmt.Printf("reply.Y %v\n", reply.Y)
	} else {
		fmt.Printf("call failed!\n")
	}
}

func register() (int32, int) {
	reply := RegisterReply{}
	call("Coordinator.Register", &RegistrArgs{}, &reply)
	return reply.Id, reply.NReduce
}

func ask(wid int32) (int32, Task) {
	args := AssignArgs{wid}
	reply := AssignReply{}
	ok := call("Coordinator.Assign", &args, &reply)
	if !ok {
		fmt.Printf("ask: call failed!\n")
	}
	return reply.Id, Task{reply.Flag, reply.File}
}

func start(wid, tid int32, flag string) {
	args := RoundArgs{wid, tid, flag}
	ok := call("Coordinator.Round", &args, &RoundReply{})
	if !ok {
		fmt.Printf("ask: call failed!\n")
	}
}

func work(tid int32, task Task, mapf mf, reducef rf) {
	switch task.Flag {
	case "map":
		log.Printf("map task: %s start\n", task.File)

		// preprocess for map
		text := process(task.File, "map").(string)
		kva := mapf(task.File, text)
		if submit(wid, tid, "map") {
			log.Printf("coordinator accept task: %s\n", task.File)
			save(kva, wid, tid, nr)
		} else {
			log.Printf("the task: %s refused by coordinator\n", task.File)
		}
	case "reduce":
		contents := make([]string, 0)
		keys := make([]string, 0)
		log.Printf("reduce task: %s start\n", task.File)

		// preprocess for reduce
		inter := process(task.File, "reduce").([]KeyValue)
		sort.Sort(ByKey(inter))
		i := 0
		for i < len(inter) {
			j := i + 1
			for j < len(inter) && inter[j].Key == inter[i].Key {
				j++
			}
			values := []string{}
			for k := i; k < j; k++ {
				values = append(values, inter[k].Value)
			}
			output := reducef(inter[i].Key, values)
			contents = append(contents, output)
			keys = append(keys, inter[i].Key)
			i = j
		}
		if submit(wid, tid, "reduce") {
			log.Printf("coordinator accept task: %s\n", task.File)
			finish(wid, tid, keys, contents)
		} else {
			log.Printf("the task: %s refused by coordinator\n", task.File)
		}
	}
}

// submit task to coordinator
func submit(wid int32, tid int32, flag string) bool {
	log.Printf("worker submit task: %d\n", tid)

	args := FinishArgs{wid, tid, flag}
	reply := &FinishReply{}
	call("Coordinator.Finish", &args, &reply)

	return reply.Accept
}

func done(tid int32) {
	args := ReviewArgs{tid}
	call("Coordinator.Review", &args, &ReviewReply{})
}

// send an RPC request to the coordinator, wait for the response.
// usually returns true.
// returns false if something goes wrong.
func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
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

func process(file string, flag string) any {
	if flag == "map" {
		f, _ := os.Open(file)
		contents, _ := io.ReadAll(f)
		return string(contents)
	} else {
		inter := []KeyValue{}
		f, _ := os.Open(file)
		b, _ := io.ReadAll(f)
		contents := strings.Split(string(b), "\n")
		for _, word := range contents {
			inter = append(inter, KeyValue{word, "1"})
		}
		return inter
	}
}

func save(kva []KeyValue, wid int32, tid int32, nr int) {
	format := "mr-out-%d-%d"
	buckets := make([][]string, nr)
	for i := range buckets {
		buckets[i] = make([]string, 0)
	}
	for _, kv := range kva {
		b := ihash(kv.Key) % nr
		buckets[b] = append(buckets[b], kv.Key)
	}
	for i, b := range buckets {
		file := fmt.Sprintf(format, tid, i)
		f, _ := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0666)
		for _, word := range b {
			fmt.Fprintf(f, "%s\n", word)
		}
		f.Close()
	}
	log.Printf("worker %d has finished map task: %d\n", wid, tid)
}

func finish(wid, tid int32, keys []string, contents []string) {
	file := "mr-out-" + strconv.Itoa(int(tid)-1) + ".result"
	f, _ := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()
	for i, key := range keys {
		fmt.Fprintf(f, "%v %v\n", key, contents[i])
	}
	log.Printf("worker %d has finished reduce task: %d\n", wid, tid)
}
