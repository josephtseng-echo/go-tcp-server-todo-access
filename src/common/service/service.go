package service

import (
	"net"
	"strconv"
)

type AcceptArgs  struct {
	Service *Service
	TcpConn *net.TCPConn
	Fd uint64
	Error chan int
	ClientIp string
	Packet chan []byte      //默认的数据包
	ICPacket chan []byte    //IC数据包
	ESPacket chan []byte    //ES数据包
	QEPacket chan []byte    //QE数据包
	RPCPacket chan []byte   //RPC数据包
	Ext Ext
}

type Service struct {
	ip       string
	port     int
	onAccept func(*AcceptArgs)
	tcpListener *net.TCPListener
}

func ListenTCP(ip string, port int) (TcpServer, error) {
	address := ip
	address += ":"
	address += strconv.Itoa(port)
	var err error
	var addr *net.TCPAddr
	var listener *net.TCPListener
	
	addr, err = net.ResolveTCPAddr("tcp", address)
	if nil != err {
		return nil, err
	}
	listener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, err
	}
	server := NewTcpServer(ip, port, listener)
	return server, nil
}