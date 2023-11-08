package raft

//
// this is an outline of the API that raft must expose to
// the service (or tester). see comments below for
// each of these functions for more details.
//
// rf = Make(...)
//   create a new Raft server.
// rf.Start(command interface{}) (index, term, isleader)
//   start agreement on a new log entry
// rf.GetState() (term, isLeader)
//   ask a Raft for its current term, and whether it thinks it is leader
// ApplyMsg
//   each time a new entry is committed to the log, each Raft peer
//   should send an ApplyMsg to the service (or tester)
//   in the same server.
//

import (
	"log"
	//	"bytes"
	"math/rand"
	"raft/labrpc"
	"sync"
	"sync/atomic"
	"time"
	//	"6.5840/labgob"
)

// as each Raft peer becomes aware that successive log entries are
// committed, the peer should send an ApplyMsg to the service (or
// tester) on the same server, via the applyCh passed to Make(). set
// CommandValid to true to indicate that the ApplyMsg contains a newly
// committed log entry.
//
// in part 2D you'll want to send other kinds of messages (e.g.,
// snapshots) on the applyCh, but set CommandValid to false for these
// other uses.
type ApplyMsg struct {
	CommandValid bool
	Command      interface{}
	CommandIndex int

	// For 2D:
	SnapshotValid bool
	Snapshot      []byte
	SnapshotTerm  int
	SnapshotIndex int
}

const (
	LEADER = iota
	FOLLOWER
	CANDIDATE
)

var (
	HeartBeatTimeout = 300 + rand.Intn(150)
	ElectionTimeout  = 200 + rand.Intn(50)
	HeartBeatPeriod  = 100
)

// A Go object implementing a single Raft peer.
type Raft struct {
	mu        sync.Mutex          // Lock to protect shared access to this peer's state
	peers     []*labrpc.ClientEnd // RPC end points of all peers
	persister *Persister          // Object to hold this peer's persisted state
	me        int                 // this peer's index into peers[]
	dead      int32               // set by Kill()

	// Your data here (2A, 2B, 2C).
	timer1 *time.Ticker // 心跳计时器
	timer2 *time.Ticker // 选举计时器

	currentTerm int32 // 当前任期
	votedFor    int32 // 投票目标
	state       int32

	// Look at the paper's Figure 2 for a description of what
	// state a Raft server must maintain.

}

// return currentTerm and whether this server
// believes it is the leader.
func (rf *Raft) GetState() (int, bool) {

	var term int
	var isleader bool
	// Your code here (2A).

	term = int(rf.getTerm())
	isleader = rf.isLeader()
	return term, isleader
}

func (rf *Raft) isLeader() bool {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	return rf.state == LEADER
}

func (rf *Raft) getTerm() int32 {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	return rf.currentTerm
}

func (rf *Raft) getState() int32 {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	return rf.state
}

// save Raft's persistent state to stable storage,
// where it can later be retrieved after a crash and restart.
// see paper's Figure 2 for a description of what should be persistent.
// before you've implemented snapshots, you should pass nil as the
// second argument to persister.Save().
// after you've implemented snapshots, pass the current snapshot
// (or nil if there's not yet a snapshot).
func (rf *Raft) persist() {
	// Your code here (2C).
	// Example:
	// w := new(bytes.Buffer)
	// e := labgob.NewEncoder(w)
	// e.Encode(rf.xxx)
	// e.Encode(rf.yyy)
	// raftstate := w.Bytes()
	// rf.persister.Save(raftstate, nil)
}

// restore previously persisted state.
func (rf *Raft) readPersist(data []byte) {
	if data == nil || len(data) < 1 { // bootstrap without any state?
		return
	}
	// Your code here (2C).
	// Example:
	// r := bytes.NewBuffer(data)
	// d := labgob.NewDecoder(r)
	// var xxx
	// var yyy
	// if d.Decode(&xxx) != nil ||
	//    d.Decode(&yyy) != nil {
	//   error...
	// } else {
	//   rf.xxx = xxx
	//   rf.yyy = yyy
	// }
}

// the service says it has created a snapshot that has
// all info up to and including index. this means the
// service no longer needs the log through (and including)
// that index. Raft should now trim its log as much as possible.
func (rf *Raft) Snapshot(index int, snapshot []byte) {
	// Your code here (2D).

}

// example RequestVote RPC arguments structure.
// field names must start with capital letters!
type RequestVoteArgs struct {
	// Your data here (2A, 2B).
	Term        int32
	CandidateId int32
}

// example RequestVote RPC reply structure.
// field names must start with capital letters!
type RequestVoteReply struct {
	// Your data here (2A).
	Term         int32
	VotedGranted bool
}

func (rf *Raft) setTerm(term int32) {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	rf.currentTerm = term
}

func (rf *Raft) incTerm() {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	rf.currentTerm++
}

// example RequestVote RPC handler.
func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	log.Printf("[%d] recieve vote request from [%d]\n", rf.me, args.CandidateId)

	// 保护任期的读写
	rf.mu.Lock()
	defer rf.mu.Unlock()

	// 处理投票请求
	// 请求任期小于当前任期
	if args.Term < rf.currentTerm {
		reply.VotedGranted = false  // 废弃该请求
		reply.Term = rf.currentTerm // 回复任期
		log.Printf("[%d] reject! [%d] earlier term(%d)\n", rf.me, args.CandidateId, args.Term)

		return
	}

	// 请求任期大于当前任期
	if args.Term > rf.currentTerm {
		log.Printf("[%d] discover later term(%d) from [%d] in term(%d). become follower!\n",
			rf.me, args.Term, args.CandidateId, rf.currentTerm)
		reply.VotedGranted = true      // 废弃该请求
		rf.currentTerm = args.Term     // 更新当前任期
		reply.Term = rf.currentTerm    // 回复任期
		rf.votedFor = args.CandidateId // 投票
		rf.state = FOLLOWER            // 转变为跟随者
		log.Printf("[%d] vote for [%d]\n", rf.me, args.CandidateId)

		return
	}

	// 接受投票请求
	if rf.votedFor == -1 || rf.votedFor == args.CandidateId {
		reply.VotedGranted = true
		reply.Term = rf.currentTerm    // 回复任期
		rf.votedFor = args.CandidateId // 投票
		rf.resetHeatBeat()
		log.Printf("[%d] vote for [%d]\n", rf.me, args.CandidateId)

		return
	}

	reply.VotedGranted = false
}

type AppendEntriesArgs struct {
	Term     int32 // leader's term
	LeaderId int32 // so follower can redirect clients
}

type AppendEntriesReply struct {
	Term    int32 // current term
	Success bool
}

// 发起心跳
func (rf *Raft) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	if args.Term < rf.currentTerm {
		reply.Success = false
		return
	}

	if args.Term > rf.currentTerm {
		reply.Success = true
		log.Printf("[%d] discover new leader of term(%d) in term(%d)\n", rf.me, args.Term, rf.currentTerm)
		rf.currentTerm = args.Term
		reply.Term = rf.currentTerm
		rf.state = FOLLOWER
		rf.votedFor = -1
		rf.resetHeatBeat()

		return
	}

	rf.state = FOLLOWER
	rf.resetHeatBeat()
	log.Printf("[%d] recieve hearbeat\n", rf.me)

	reply.Success = true
}

// example code to send a RequestVote RPC to a server.
// server is the index of the target server in rf.peers[].
// expects RPC arguments in args.
// fills in *reply with RPC reply, so caller should
// pass &reply.
// the types of the args and reply passed to Call() must be
// the same as the types of the arguments declared in the
// handler function (including whether they are pointers).
//
// The labrpc package simulates a lossy network, in which servers
// may be unreachable, and in which requests and replies may be lost.
// Call() sends a request and waits for a reply. If a reply arrives
// within a timeout interval, Call() returns true; otherwise
// Call() returns false. Thus Call() may not return for a while.
// A false return can be caused by a dead server, a live server that
// can't be reached, a lost request, or a lost reply.
//
// Call() is guaranteed to return (perhaps after a delay) *except* if the
// handler function on the server side does not return.  Thus there
// is no need to implement your own timeouts around Call().
//
// look at the comments in ../labrpc/labrpc.go for more details.
//
// if you're having trouble getting RPC to work, check that you've
// capitalized all field names in structs passed over RPC, and
// that the caller passes the address of the reply struct with &, not
// the struct itself.
func (rf *Raft) sendRequestVote(server int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	return ok
}

// the service using Raft (e.g. a k/v server) wants to start
// agreement on the next command to be appended to Raft's log. if this
// server isn't the leader, returns false. otherwise start the
// agreement and return immediately. there is no guarantee that this
// command will ever be committed to the Raft log, since the leader
// may fail or lose an election. even if the Raft instance has been killed,
// this function should return gracefully.
//
// the first return value is the index that the command will appear at
// if it's ever committed. the second return value is the current
// term. the third return value is true if this server believes it is
// the leader.
func (rf *Raft) Start(command interface{}) (int, int, bool) {
	index := -1
	term := -1
	isLeader := true

	// Your code here (2B).

	return index, term, isLeader
}

// the tester doesn't halt goroutines created by Raft after each test,
// but it does call the Kill() method. your code can use killed() to
// check whether Kill() has been called. the use of atomic avoids the
// need for a lock.
//
// the issue is that long-running goroutines use memory and may chew
// up CPU time, perhaps causing later tests to fail and generating
// confusing debug output. any goroutine with a long-running loop
// should call killed() to check whether it should stop.
func (rf *Raft) Kill() {
	atomic.StoreInt32(&rf.dead, 1)
	// Your code here, if desired.
}

func (rf *Raft) killed() bool {
	z := atomic.LoadInt32(&rf.dead)
	return z == 1
}

func (rf *Raft) ticker() {
	for rf.killed() == false {

		// Your code here (2A)
		// Check if a leader election should be started.

		var wg sync.WaitGroup
		for peer := range rf.peers {
			if rf.getState() != LEADER {
				break
			}

			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				rf.heartBeat(i)
			}(peer)
		}

		// pause for a random amount of time between 50 and 350
		// milliseconds.
		ms := HeartBeatPeriod + rand.Intn(5)
		time.Sleep(time.Duration(ms) * time.Millisecond)
	}
}

// the service or tester wants to create a Raft server. the ports
// of all the Raft servers (including this one) are in peers[]. this
// server's port is peers[me]. all the servers' peers[] arrays
// have the same order. persister is a place for this server to
// save its persistent state, and also initially holds the most
// recent saved state, if any. applyCh is a channel on which the
// tester or service expects Raft to send ApplyMsg messages.
// Make() must return quickly, so it should start goroutines
// for any long-running work.
func Make(peers []*labrpc.ClientEnd, me int,
	persister *Persister, applyCh chan ApplyMsg) *Raft {
	rf := &Raft{}
	rf.peers = peers
	rf.persister = persister
	rf.me = me
	HeartBeatTimeout = 150 + rand.Intn(150)
	rf.timer1 = time.NewTicker(time.Duration(HeartBeatTimeout))
	rf.timer2 = time.NewTicker(time.Duration(ElectionTimeout))

	// Your initialization code here (2A, 2B, 2C).
	rf.state = FOLLOWER

	go rf.heartBeatTimeOut()
	// initialize from state persisted before a crash
	rf.readPersist(persister.ReadRaftState())

	// start ticker goroutine to start elections
	go rf.ticker()

	return rf
}

// 重置心跳计时
func (rf *Raft) resetHeatBeat() {
	// 设置随机种子，以确保每次运行生成的随机数不同
	rand.Seed(time.Now().UnixNano())

	// 生成随机时间间隔
	d := time.Duration(300+rand.Intn(150)) * time.Millisecond
	rf.timer1.Reset(d)
}

// 重置选举计时
func (rf *Raft) resetElection() {
	// 设置随机种子，以确保每次运行生成的随机数不同
	rand.Seed(time.Now().UnixNano())
	d := time.Duration(ElectionTimeout) * time.Millisecond
	rf.timer2.Reset(d)
}

func (rf *Raft) heartBeatTimeOut() {
	for range rf.timer1.C {
		if rf.getState() != LEADER {
			log.Printf("[%d] heartbeat timeout!", rf.me)
			rf.election()
		}
	}
}

func (rf *Raft) electionTimeOut() {
	for range rf.timer2.C {
		log.Printf("[%d] election timeout!", rf.me)
		rf.election()
	}
}

func (rf *Raft) setState(state int32) {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	rf.state = state
}

// 开启选举
func (rf *Raft) election() {
	rf.resetElection() // 重置选举计时器
	var voted int32 = 0

	rf.mu.Lock()
	rf.state = CANDIDATE
	rf.currentTerm++ // 任期递增
	rf.votedFor = -1
	rf.mu.Unlock()

	log.Printf("[%d] start election of term(%d)\n", rf.me, rf.currentTerm)

	var wg sync.WaitGroup
	voteCh := make(chan struct{}, len(rf.peers))

	args := RequestVoteArgs{
		Term:        rf.currentTerm,
		CandidateId: int32(rf.me),
	}

	go func(ch chan struct{}) {
		for range ch {
			voted++

			if rf.getState() != CANDIDATE {
				return
			}
			// 获票数超过半数 当选本轮领导
			if int(voted) > len(rf.peers)/2 {
				rf.mu.Lock()
				if rf.state == CANDIDATE {
					rf.state = LEADER
					log.Printf("[%d] become the leader of term(%d)!\n", rf.me, rf.currentTerm)
					rf.resetHeatBeat()
				}
				rf.mu.Unlock()
				return
			}
		}
	}(voteCh)

	for i := 0; i < len(rf.peers); i++ {
		// 已经不是候选者
		if rf.getState() != CANDIDATE {
			return
		}

		wg.Add(1)
		go func(peer int) {
			defer wg.Done()

			log.Printf("[%d] send vote request to [%d]\n", rf.me, peer)
			// 发送投票请求
			reply := RequestVoteReply{}
			if !rf.sendRequestVote(peer, &args, &reply) {
				log.Printf("[%d] send vote request to [%d] failed\n", rf.me, peer)
				return
			}

			// 请求投票成功
			if reply.VotedGranted {
				voteCh <- struct{}{}
				log.Printf("[%d] voted by [%d]! total voted: %d\n", rf.me, peer, voted)
				return
			}

			// 回复任期大于当前任期
			rf.mu.Lock()
			if reply.Term > rf.currentTerm {
				log.Printf("[%d] discover later term(%d) from [%d] in term(%d). become follower!\n",
					rf.me, reply.Term, peer, rf.currentTerm)
				rf.currentTerm = reply.Term
				rf.state = FOLLOWER
			}
			rf.mu.Unlock()
		}(i)
	}

	wg.Wait()
}

// 发起心跳
func (rf *Raft) heartBeat(server int) bool {
	if rf.me == server {
		return false
	}

	log.Printf("[%d] send heartbeat to [%d]\n", rf.me, server)

	args := AppendEntriesArgs{
		Term:     rf.getTerm(),
		LeaderId: int32(rf.me),
	}

	reply := AppendEntriesReply{}

	// 发送心跳
	ok := rf.peers[server].Call("Raft.AppendEntries", &args, &reply)
	// 发现回复任期大于当前任期
	rf.mu.Lock()
	if reply.Term > rf.currentTerm {
		rf.currentTerm = reply.Term // 更新任期
		rf.state = FOLLOWER
	}
	rf.mu.Unlock()

	return ok
}
