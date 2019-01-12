package service

import (
	"common/base/global"
	"net"
)

type TcpServer interface {
	Start()
	Stop()
	OnAccept(f func(*AcceptArgs))
	OnListener()(*net.TCPListener)
}

func (service *Service) OnAccept(f func(*AcceptArgs)) {
	service.onAccept = f
}

func (service *Service) Stop() {
	_, ok := <-global.ServiceCh
	if ok {
		log := global.Config.Logger
		log.Info("service stop.")
		close(global.ServiceCh)
		return
	}
	global.ServiceWg.Wait()
}

func (service *Service) Start() {
	defer global.ServiceWg.Done()
	log := global.Config.Logger
	log.Info("(local) service start.")
	var err error
	var conn *net.TCPConn
	listener := service.OnListener()
	for {
		select {
			case _, ok := <-global.ServiceCh:
				if ok {
					return
				}
			default:
		}
		conn, err = listener.AcceptTCP()
		if err != nil {
			return
		}
		f, _ := conn.File()
		fd := f.Fd()
		args := &AcceptArgs {
			ClientIp : conn.RemoteAddr().String(),
			Service : service,
			TcpConn : conn,
			Fd : uint64(fd),	
		}
		global.ServiceWg.Add(1)
		go func(args *AcceptArgs) {
			defer global.ServiceWg.Done()
			service.onAccept(args)
		}(args)
	}
}

func (service *Service) OnListener()(*net.TCPListener) {
	return service.tcpListener
}

func NewTcpServer(ip string, port int, listener *net.TCPListener) (TcpServer) {
	return &Service{
		ip:   ip,
		port: port,
		tcpListener: listener,
	}
}