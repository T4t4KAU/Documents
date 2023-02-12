package mr

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	READY = iota
	MAP
	SUSPEND
	REDUCE
	DONE
)

type state struct {
	worker int32
	done   chan struct{}
}

type Task struct {
	Flag string
	File string
}

type Coordinator struct {
	states  map[int32]*state // tid -> state
	tasks   map[int32]Task   // tid -> task
	stage   int              // map or reduce
	nmap    int              // num of map tasks
	nreduce int              // num of reduce tasks
	queue   *queue           // task queue
	nworker int32            // num of workers
	buckets []string         // reduce buckets
	mux     sync.Mutex
	finish  chan struct{}
}

// the RPC argument and reply types are defined in rpc.go.
func (c *Coordinator) Example(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = args.X + 1
	return nil
}

// assign worker id
func (c *Coordinator) Register(args *RegistrArgs, reply *RegisterReply) error {
	if c.stage == DONE {
		return errors.New("coordinator: tasks all done")
	}

	wid := atomic.AddInt32(&c.nworker, 1) // generate worker id
	reply.Id, reply.NReduce = wid, int(c.nreduce)
	log.Printf("assign worker id: %d\n", wid)
	return nil
}

// assign task for worker
func (c *Coordinator) Assign(args *AssignArgs, reply *AssignReply) error {
	if args.Id <= 0 {
		return errors.New("coordinator: wrong worker id")
	}

	if c.stage == SUSPEND || c.stage == DONE || c.stage == READY {
		reply.Id = int32(-c.stage)
		return nil
	}

	if !c.queue.empty() {
		item := c.queue.Dequeue() // get map task from queue
		reply.Id, reply.File = item.task, item.file
		switch c.stage {
		case MAP:
			reply.Flag = "map"
		case REDUCE:
			reply.Flag = "reduce"
		}
		log.Printf("assign %s task: %d for worker %d\n",
			reply.Flag, reply.Id, args.Id)
	}

	return nil
}

func (c *Coordinator) Round(args *RoundArgs, reply *RoundReply) error {
	task, ok := c.tasks[args.TID]
	if !ok || task.Flag != args.Flag {
		return errors.New("coordinator: nonexistent task")
	}

	log.Printf("worker %d start %s task %d\n", args.WID, args.Flag, args.TID)
	done := make(chan struct{})
	c.mux.Lock()
	c.states[args.TID] = &state{args.WID, done}
	c.mux.Unlock()

	go func() {
		select {
		case <-done:
			return
		// task timeout
		case <-time.After(10 * time.Second):
			log.Printf("task %d of worker %d timeout\n", args.TID, args.WID)
			c.mux.Lock()
			delete(c.states, args.TID) // remove
			c.mux.Unlock()
			c.queue.Enqueue(item{args.TID, task.File}) // recover queue
		}
	}()

	return nil
}

func (c *Coordinator) Finish(args *FinishArgs, reply *FinishReply) error {
	st, ok := c.states[args.TID]
	if !ok || st.worker != args.WID {
		reply.Accept = false
		return errors.New("coordinator: invalid worker id")
	}
	if args.Flag == "map" && c.stage == REDUCE ||
		c.stage == SUSPEND || c.stage == DONE {
		reply.Accept = false
		return errors.New("coordinator: invalid task")
	}

	reply.Accept = true
	st.done <- struct{}{}

	c.mux.Lock()
	delete(c.states, args.TID)
	c.mux.Unlock()

	return nil
}

func (c *Coordinator) Review(args *ReviewArgs, reply *ReviewReply) error {
	if c.check(args.Id) {
		log.Printf("worker %d has finished task %d\n", args.Id, args.Id)
		c.finish <- struct{}{}
		return nil
	}
	log.Printf("worker %d has not finished task %d\n", args.Id, args.Id)
	return errors.New("coordinator: review failed")
}

// commit files
func (c *Coordinator) commit(tid int, done chan struct{}) {
	for b := 0; b < c.nreduce; b++ {
		file := fmt.Sprintf("mr-out-%d-%d", tid, b)
		log.Printf("commit %s...\n", file)
		f, err := os.Open(file)
		if err != nil {
			log.Printf("open %s error: %v\n", file, err)
			// done <- struct{}{}
			return
		}
		bytes, err := io.ReadAll(f)
		if err != nil {
			log.Printf("read %s error: %v\n", file, err)
		}
		f.Close()
		os.Remove(file)

		// c.mux.Lock()
		c.buckets[b] += string(bytes)
		// c.mux.Unlock()
	}
	// done <- struct{}{}
}

// merge output of map tasks
// change to reduce stage :)
func (c *Coordinator) merge() {
	done := make(chan struct{}, c.nmap)

	log.Println("merge map output...")
	for i := 0; i < c.nmap; i++ {
		c.commit(i+1, done)
	}
	// for i := 0; i < c.nmap; i++ {
	// 	<-done
	// }
	for i := 0; i < c.nreduce; i++ {
		file := "mr-out-" + strconv.Itoa(i)
		log.Printf("create reduce file: %s\n", file)
		f, _ := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0666)
		f.WriteString(strings.Trim(c.buckets[i], "\n"))
		f.Close()
	}
}

func (c *Coordinator) check(tid int32) bool {
	format := "mr-out-%d"
	if c.stage == MAP {
		log.Printf("check map task %d...\n", tid)
		format += "-%d"
		for i := 0; i < c.nreduce; i++ {
			_, err := os.Stat(fmt.Sprintf(format, tid, i))
			if os.IsNotExist(err) {
				return false
			}
		}
	} else if c.stage == REDUCE {
		log.Printf("check reduce task %d...\n", tid)
		format += ".result"
		_, err := os.Stat(fmt.Sprintf(format, tid-1))
		if os.IsNotExist(err) {
			return false
		}
	} else {
		return false
	}
	return true
}

// start a thread that listens for RPCs from worker.go
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
func (c *Coordinator) Done() bool {
	return c.stage == DONE
}

// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{}

	c.states = make(map[int32]*state)
	c.tasks = make(map[int32]Task)
	c.finish = make(chan struct{}, c.nmap)
	c.buckets = make([]string, nReduce)
	c.queue = Queue()
	c.nmap = len(files)
	c.nreduce = nReduce
	c.stage = READY

	// init map and queue
	for i, file := range files {
		c.tasks[int32(i+1)] = Task{"map", file}
		c.queue.Enqueue(item{int32(i + 1), file})
	}

	c.stage = MAP

	go func() {
		for i := 0; i < c.nmap; i++ {
			<-c.finish
		}

		log.Printf("map tasks all done\n")
		if !c.queue.empty() {
			log.Panicln("queue is not empty")
		}

		c.stage = SUSPEND
		c.merge()
		log.Printf("change to reduce stage\n")
		for i := 0; i < c.nreduce; i++ {
			file := "mr-out-" + strconv.Itoa(i)
			log.Printf("%s enter the queue\n", file)
			c.queue.Enqueue(item{int32(i + 1), file})
			c.tasks[int32(i+1)] = Task{"reduce", file}
		}

		c.stage = REDUCE

		for i := 0; i < c.nreduce; i++ {
			<-c.finish
		}

		for i := 0; i < c.nreduce; i++ {
			os.Remove("mr-out-" + strconv.Itoa(i))
		}

		c.stage = DONE
	}()

	c.server()
	return &c
}
