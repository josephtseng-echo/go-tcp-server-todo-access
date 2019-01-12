package service

import (
	"common/base/global"
	"net"
	"time"
)

type EasyClientInterface interface {
	SendNotRecvData([]byte) error
	SendAndRecvData([]byte) ([]byte, error)
	GetClientAddr() string
}

type EasyClient struct {
	Conn       net.Conn
	ClientAddr string
}

func NewTcpClient(clientAddr string) EasyClientInterface {
	conn, err := net.Dial("tcp", clientAddr)
	if err != nil {
		global.Config.Logger.Error("(local) client new tcp ", clientAddr, " error: ", err.Error())
	} else {
		global.Config.Logger.Info("(local) client new tcp ", clientAddr, " ok ")
	}
	return &EasyClient{
		Conn:       conn,
		ClientAddr: clientAddr,
	}
}

func (c EasyClient) GetClientAddr() string {
	return c.ClientAddr
}

func (c EasyClient) SendNotRecvData(send []byte) error {
	var err error
	_, err = c.Conn.Write(send)
	if err != nil {
		return err
	}
	return nil
}

func (c EasyClient) SendAndRecvData(send []byte) ([]byte, error) {
	var n int
	var err error
	var readBuf = make([]byte, 2048)
	n, err = c.Conn.Write(send)
	if err != nil {
		return nil, err
	}
	c.Conn.SetReadDeadline(time.Now().Add(time.Duration(5) * time.Second))
	n, err = c.Conn.Read(readBuf)
	if err != nil {
		return nil, err
	} else {
		return readBuf[:n], nil
	}
}
