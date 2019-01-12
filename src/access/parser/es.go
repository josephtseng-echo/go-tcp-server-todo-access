package parser

import (
	"bytes"
	"common/base/function"
	"common/base/global"
	"common/base/packet"
	"common/ext"
	"common/service"
	"encoding/binary"
	"fmt"
)

func OnParserEs(args *service.AcceptArgs, buf []byte) {
	defer global.ServiceWg.Done()
	var bodyBuf []byte
	var Cmd uint16
	var ip string
	var ipNew string
	var port int32
	head := new(EsHead)
	c, err := EsHeadDecode(buf[:ESHEADSIZE], head)
	if c == false {
		global.Config.Logger.Error("es head decode error : ", err.Error())
		return
	}
	Cmd = head.Cmd
	bodyBuf = EsGetBodyPacket(buf, head.Len)
	if Cmd != 0x1010 {
		global.Config.Logger.Info("[client] send cmd = ", fmt.Sprintf("0x%x", Cmd), " to server")
		global.Config.Logger.Info("[client] head=", fmt.Sprintf("%+v", head))
	}

	if args.Fd <= 0 {
		//非法连接
		return
	}
	task := service.GetTcpTask(args.Fd)
	task.Client.Ip = args.ClientIp
	if task == nil {
		//非法连接
		return
	}

	//用户登录成功后会绑定userid
	if task.UserId <= 0 && task.Status != 2 {
		//TODO
		if Cmd == 0x6001 || Cmd == 0x1001 {
			if Cmd == 0x1001 {
				task.Ext.FirstPacket = buf
			}
			if Cmd == 0x6001 {
				//string = ip, int = port, int = version 包体
				//TODO 补充对异常处理
				reader := packet.Reader(bodyBuf)
				ip, _ = reader.ReadString()
				ipNew = function.GetValidString(ip) //一定要过滤未的０,lua客户端的一个bug
				port, _ = reader.ReadInt32()
				version, _ := reader.ReadInt32()
				task.Ext.MyServerIp = ip
				task.Ext.MyServerPort = port
				task.Ext.MyServerVersion = version
				task.Ext.ClientServer = ext.AddEsClient(ipNew, port)
				task.Ext.ClientServer = ext.AddEsClient(ipNew, port)
			}
			if len(task.Ext.FirstPacket) > 0 && Cmd == 0x1001 {
				reader := packet.Reader(EsGetBodyPacket(task.Ext.FirstPacket, head.Len))
				tid, _ := reader.ReadInt32()
				uid, _ := reader.ReadInt32()
				reader.ReadString() //mtkey
				reader.ReadString()
				//绑定task userid
				task.UserId = uint64(uid)
				//测试用的tid
				tid = 65536

				task.Status = 2
				task.Ext.MyTid = int64(tid)
				task.Ext.FirstPacket = nil
			}
		} else {
			//TODO 未登录其他包的处理
		}
	}
	//登录后
	//心跳包
	if Cmd == 0x1010 {
		//
	} 
	if Cmd != 0x1010 {
		global.Config.Logger.Info("[client] task=", fmt.Sprintf("%+v", task))
	}
}

func EsGetBodyPacket(buf []byte, length uint16) []byte {
	offset := ESHEADSIZE + int(length)
	bufNew := buf[ESHEADSIZE:offset]
	return bufNew
}

func EsHeadDecode(buf []byte, esHead *EsHead) (bool, error) {
	var err error
	headBuf := bytes.NewBuffer(buf)
	err = binary.Read(headBuf, binary.BigEndian, &esHead.Flag)
	if err != nil {
		return false, err
	}
	err = binary.Read(headBuf, binary.BigEndian, &esHead.Cmd)
	if err != nil {
		return false, err
	}
	err = binary.Read(headBuf, binary.BigEndian, &esHead.Ver)
	if err != nil {
		return false, err
	}
	err = binary.Read(headBuf, binary.BigEndian, &esHead.Len)
	if err != nil {
		return false, err
	}
	return true, nil
}

func EsHeadEncode(cmd uint16, buf []byte) []byte {
	length := len(buf)
	head := &EsHead{
		Flag: [2]byte{'E', 'S'},
		Cmd:  cmd,
		Ver:  2,
		Len:  0,
	}
	head.Len = uint16(length)
	wBuf := new(bytes.Buffer)
	binary.Write(wBuf, binary.BigEndian, &head.Flag)
	binary.Write(wBuf, binary.BigEndian, &head.Cmd)
	binary.Write(wBuf, binary.BigEndian, &head.Ver)
	binary.Write(wBuf, binary.BigEndian, &head.Len)
	wBuf.Write(buf)
	return wBuf.Bytes()
}
