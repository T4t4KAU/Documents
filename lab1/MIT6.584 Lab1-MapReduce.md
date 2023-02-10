# MIT6.584 Lab1-MapReduce

Lab地址： https://pdos.csail.mit.edu/6.824/labs/lab-mr.html

## 初步探究

按照实验要求，要实现一个Distributed MapReduce，功能是Word Count，关于MapReduce的原理在论文里了，就不再赘述了，这里着重实验部分，先熟悉一下官方给出的相关代码

按照提示，先编译wc.go，会生成一个wc.so

```powershell
$ go build -buildmode=plugin ../mrapps/wc.go
```

其中包含了Map，Reduce这两个函数

然后执行：

```powershell
$ go run mrworker.go wc.so
```

运行完毕后就会生成保存结果的文件，执行一下命令输出Word Count的结果

```powershell
$ cat mr-out-* | sort | more
```

上述是一个单机的MapReduce，现在实验要求实现一个分布式的MapReduce，其中包含了两个角色，即Coordinator和Worker，每个Worker向Coordinator请求任务，程序要基于RPC

先将与本实验相关的文件复制到单独的文件夹中，再执行`go mod init`，初始化go.mod文件，目录结构：

```powershell
.
├── go.mod
├── main
│   ├── mrcoordinator.go
│   ├── mrworker.go
│   ├── pg-being_ernest.txt
│   ├── pg-dorian_gray.txt
│   ├── pg-frankenstein.txt
│   ├── pg-grimm.txt
│   ├── pg-huckleberry_finn.txt
│   ├── pg-metamorphosis.txt
│   ├── pg-sherlock_holmes.txt
│   ├── pg-tom_sawyer.txt
│   └── wc.so
└── mr
    ├── coordinator.go
    ├── rpc.go
    └── worker.go
```

下面增加一些基本的代码，实现Worker和Coordinator的RPC通信

增加RPC调用的参数和返回值结构，用于任务的分配

```go
type DispatchArgs struct {
	Id int
}

type DispatchReply struct {
	Id   int
	Task string
	File string
}
```

当worker向coordinator请求任务时，coordinator至少应该返回一个文件名，交予worker来处理

定义Coordianator结构：

```go
type Coordinator struct {
	files   []string
	workers map[int]net.Listener
}
```

完善其初始化函数：

```go
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{}

	c.files = files
	c.workers = make(map[int]net.Listener)

	c.server()
	return &c
}
```

增加RPC调用函数

```go
func (c *Coordinator) Dispatch(args *DispatchArgs, reply *DispatchReply) error {
	rand.Seed(time.Now().UnixNano())
	idx := rand.Intn(len(c.files))
	reply.Task = c.files[idx]

	return nil
}
```

该函数会随机选择一个文件名发给worker

在worker.go中增加一个函数，用于请求任务：

```go
func AskTasks() {
	args := DispatchArgs{}
	reply := DispatchReply{}
	ok := call("Coordinator.Dispatch", &args, &reply)
	if ok {
		fmt.Printf("task: %v\n", reply.Task)
	} else {
		fmt.Printf("call failed!\n")
	}
}
```

修改Woker函数：

```go
func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {		
	AskTasks()
}
```

运行：

先启动coordinator：

```powershell
$ go run mrcoordinator.go pg-*.txt
```

再启动worker：

```powershell
$ go run mrworker.go wc.so 
task: pg-tom_sawyer.txt
```

worker成功接收到了文件名，表明成功实现RPC通信

