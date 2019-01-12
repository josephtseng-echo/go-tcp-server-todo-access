package task

import (
	"access/parser"
	"bufio"
	"bytes"
	"common/base/global"
	"common/service"
	"encoding/binary"
	"io"
)

func onError(args *service.AcceptArgs) {
	defer global.ServiceWg.Done()
	for {
		select {
		case i, ok := <-args.Error:
			if ok {
				mLog := global.Config.Logger
				mLog.Error("error: ", i)
			}
		}
	}
}

func onRead(args *service.AcceptArgs) {
	defer global.ServiceWg.Done()
	for {
		select {

		case buf, ok := <-args.ESPacket:
			if ok {
				global.ServiceWg.Add(1)
				go parser.OnParserEs(args, buf)
			}

		case buf, ok := <-args.RPCPacket:
			if ok {
				global.ServiceWg.Add(1)
				go parser.OnParserRpc(args, buf)
			}
		}
	}
}

func onPacket(buf []byte, args *service.AcceptArgs) []byte {
	var length uint16
	var headlen uint8
	var bodylen uint16
	var i uint16
	length = uint16(len(buf))
	for i = 0; i < length; i++ {
		if i+uint16(parser.RPCHEADSIZE) > length {
			if i+uint16(parser.ESHEADSIZE) > length {
				break
			}
			break
		}
		if buf[i] == uint8('E') && buf[i+1] == uint8('S') {
			headlen = uint8(parser.ESHEADSIZE)
			binary.Read(bytes.NewReader(buf[i+6:i+8]), binary.BigEndian, &bodylen)
			if i+uint16(headlen)+bodylen > length {
				break
			}
			if i+uint16(headlen)+bodylen > uint16(parser.PACKETTOTALSIZE) {
				break
			}
			//ES
			args.ESPacket <- buf[i : i+uint16(headlen)+bodylen]
			i += uint16(headlen) + bodylen - 1
		} else if buf[i] == uint8('R') && buf[i+1] == uint8('P') && buf[i+2] == uint8('C') {
			//TODO
		}
	}
	if i == length {
		return make([]byte, 0)
	}
	return buf[i:]
}

func OnAccept(args *service.AcceptArgs) {
	mLog := global.Config.Logger
	service.AddTcpTask(args)
	reader := bufio.NewReader(args.TcpConn)
	var buf = make([]byte, 0)
	var err error
	var n int
	var tempBuf = make([]byte, parser.BUFREADTOTALSIZE)
	//默认 暂时可不需要
	//args.Packet = make(chan []byte, PACKETCHANSIZE)
	//es
	args.ESPacket = make(chan []byte, parser.PACKETCHANSIZE)
	//rpc
	args.RPCPacket = make(chan []byte, parser.PACKETCHANSIZE)

	args.Error = make(chan int, parser.ERRORCHANSIZE)
	global.ServiceWg.Add(2)
	go onRead(args)
	go onError(args)
	for {
		select {
		case _, ok := <-global.ServiceCh:
			if ok {
				return
			}
		default:
		}
		n, err = reader.Read(tempBuf)
		if err == io.EOF {
			mLog.Info("fd = ", args.Fd, args.Service, args.TcpConn, " client close")
			service.RemoveTcpTask(args.Fd)
			return
		}
		if err != nil {
			return
		}
		buf = onPacket(append(buf, tempBuf[:n]...), args)
	}
}
