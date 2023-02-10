package mr

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type task struct {
	flag string
	file string
}

type state struct {
	id     int
	worker int
	done   chan struct{}
}

type Coordinator struct {
	files    []string
	tasks    map[int32]task
	states   map[task]*state
	maps     *queue
	reduces  *queue
	nworker  int32
	ntask    int32
	finished chan struct{}
	mux      sync.RWMutex
}

// RPC handlers for the worker to call.
func (c *Coordinator) Dispatch(args *DispatchArgs, reply *DispatchReply) error {
	if args.Id <= 0 {
		return errors.New("wrong worker id")
	}

	if !c.maps.isEmpty() {
		reply.Flag = "map"
		reply.File = c.maps.Dequeue()
		fmt.Printf("coordinator dispatch map task:"+
			"%s for worker %d\n", reply.File, args.Id)

		goto END
	}

	if !c.reduces.isEmpty() {
		reply.Flag = "reduce"
		reply.File = c.reduces.Dequeue()
		fmt.Printf("coordinator dispatch reduce task:"+
			"%s for worker %d\n", reply.File, args.Id)

		goto END
	}

END:
	done := make(chan struct{})
	t := task{reply.Flag, reply.File}
	id := atomic.AddInt32(&c.ntask, 1)
	c.mux.Lock()
	defer c.mux.Unlock()
	c.states[t] = &state{int(id), args.Id, done}

	return nil
}

func (c *Coordinator) Register(args *RegisterArgs, reply *RegisterReply) error {
	atomic.AddInt32(&c.nworker, 1)
	reply.Id = int(c.nworker)
	return nil
}

func (c *Coordinator) Run(args *RunArgs, reply *RunReply) error {
	t := task{args.Flag, args.File}
	st, ok := c.states[t]
	if !ok {
		return errors.New("nonexistent task")
	}
	if st.worker != args.Id {
		return errors.New("wrong worker id: " + strconv.Itoa(args.Id))
	}

	c.mux.Lock()
	c.tasks[int32(st.id)] = t
	c.mux.Unlock()

	go func() {
		select {
		case <-st.done:
			return
		case <-time.After(10 * time.Second):
			switch args.Flag {
			case "map":
				c.maps.Enqueue(args.File)
			case "reduce":
				c.reduces.Enqueue(args.File)
			}
			c.mux.Lock()
			delete(c.states, task{args.Flag, args.File})
			c.mux.Unlock()
		}
	}()

	return nil
}

func (c *Coordinator) Finish(args *FinishArgs, reply *FinishReply) error {
	t := c.tasks[int32(args.Task)]
	st, ok := c.states[t]
	if !ok {
		return errors.New("nonexistent task")
	}
	if st.worker != args.Id {
		return errors.New("wrong worker id")
	}

	c.mux.Lock()
	delete(c.states, t)
	c.mux.Unlock()

	fmt.Printf("worker %d has finished %s task: %s\n",
		args.Id, t.flag, t.file)

	if t.flag == "map" {
		name := "mr-out-" + strconv.Itoa(args.Task)
		c.reduces.Enqueue(name)
	}

	if t.flag == "reduce" {
		c.finished <- struct{}{}
	}

	return nil
}

// start a thread that listens for RPCs from worker.go
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()
	os.Remove(sockname)
	lis, err := net.Listen("unix", sockname)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	go http.Serve(lis, nil)
}

// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
func (c *Coordinator) Done() bool {
	ret := false
	c.Wait()
	return ret
}

func (c *Coordinator) Wait() {
	for i := 0; i < len(c.files); i++ {
		<-c.finished
	}
}

// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{}

	c.files = files
	c.tasks = make(map[int32]task)
	c.states = make(map[task]*state)
	c.maps = Queue()
	c.reduces = Queue()
	c.finished = make(chan struct{})
	c.mux = sync.RWMutex{}

	for _, f := range c.files {
		c.maps.Enqueue(f)
	}

	c.server()
	return &c
}
