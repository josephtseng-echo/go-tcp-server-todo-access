package service

import (
	"net"
)

//存客户端信息
type Client struct {
	Ip string
}

//用于扩展是否使用
type Ext struct {
	MyServerIp string
	MyServerPort int32
	MyServerVersion int32
	FirstPacket []byte
	ClientServer EasyClientInterface
	MySvid int64 //用户server id
	MyTid  int64 //用户 tid
}

type Task struct {
	Fd uint64
	UserId uint64
	ServId uint32
	TcpConn *net.TCPConn
	Service *Service
	IsServ uint8
	Client Client
	Ext Ext
	Status uint8 //用于定时器清理非状态的用户连接
}

func AddTcpTask(args *AcceptArgs) *Task {
	defer Manager.M.Unlock()
	Manager.M.Lock()
	ret, ok := Manager.Tasks[args.Fd]
	if ok {
		return ret
	}
	task := &Task{
		Fd : args.Fd,
		TcpConn : args.TcpConn,
		Service : args.Service,
	}
	Manager.Tasks[args.Fd] = task
	return task
}

func GetTcpTask(fd uint64) *Task {
	defer Manager.M.RUnlock()
	Manager.M.RLock()
	ret, ok := Manager.Tasks[fd]
	if ok {
		return ret
	}
	return nil
}

func GetTcpTaskByUid(uid uint64) *Task {
	defer Manager.M.RUnlock()
	Manager.M.RLock()
	for _, tsk := range Manager.Tasks {
		if tsk.UserId == uid {
			return tsk
		}
	}
	return nil
}

func RemoveTcpTask(fd uint64) {
	_, ok := Manager.Tasks[fd]
	if ok {
		Manager.M.Lock()
		Manager.Tasks[fd].TcpConn.Close()
		delete(Manager.Tasks, fd)
		Manager.M.Unlock()
	}
}

func RemoveTcpTaskByUid(uid uint64) {
	for k, tsk := range Manager.Tasks {
		if tsk.UserId == uid {
			Manager.M.Lock()
			Manager.Tasks[k].TcpConn.Close()
			delete(Manager.Tasks, k)
			Manager.M.Unlock()
		}
	}
}