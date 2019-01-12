package intertask

import (
	"bufio"
	"bytes"
	"common/base/global"
	"common/service"
	"encoding/binary"
	"fmt"
	"io"
	"access/parser"
)

func onError(args *service.AcceptArgs) {
	defer global.ServiceWg.Done()
	for {
		select {
		case i, ok := <-args.Error:
			if ok {
				fmt.Printf("error: %i\n", i)
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
				go parser.OnParserAdmin(args, buf)
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
		if i+uint16(parser.ESHEADSIZE) > length {
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
		} 
	}
	if i == length {
		return make([]byte, 0)
	}
	return buf[i:]
}

func OnAccept(args *service.AcceptArgs) {
	service.AddTcpTask(args)
	reader := bufio.NewReader(args.TcpConn)
	var buf = make([]byte, 0)
	var err error
	var n int
	var tempBuf = make([]byte, parser.BUFREADTOTALSIZE)
	args.ESPacket = make(chan []byte, parser.PACKETCHANSIZE)
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
			service.RemoveTcpTask(args.Fd)
			return
		}
		if err != nil {
			return
		}
		buf = onPacket(append(buf, tempBuf[:n]...), args)
	}
}
