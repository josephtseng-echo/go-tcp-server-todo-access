//
// main.go
// Copyright (C) 2018 josephzeng <josephzeng36@gmail.com>
//
// Distributed under terms of the MIT license.
//

package main

import (
	"access/intertask"
	"access/task"
	"common/base/global"
	"common/service"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	initArgs()
	c, cErr := global.InitConfiger()
	if c == false {
		fmt.Println("配置有问题 ERROR:" + cErr.Error())
		os.Exit(1)
	}
	global.InitLogger("access")
	//
	var err error
	var server service.TcpServer
	mLog := global.Config.Logger
	//server init
	//外用提供给网部使用
	serverIp := global.Config.Configer.Section("server").Key("ip").Value()
	serverPortStr := global.Config.Configer.Section("server").Key("port").Value()
	serverPort, _ := strconv.Atoi(serverPortStr)
	server, err = service.ListenTCP(serverIp, serverPort)
	if err != nil {
		log.Fatal(err.Error())
	}
	server.OnAccept(task.OnAccept)
	global.ServiceWg.Add(1)
	go server.Start()
	mLog.Info("(local) out server > ", serverIp, ":", serverPort, " start success.")
	
	//新增一个端口监听　主要作用内网，管理服务使用
	interServerIp := global.Config.Configer.Section("interserver").Key("ip").Value()
	interServerPortStr := global.Config.Configer.Section("interserver").Key("port").Value()
	interServerPort, _ := strconv.Atoi(interServerPortStr)
	server, err = service.ListenTCP(interServerIp, interServerPort)
	server.OnAccept(intertask.OnAccept)
	global.ServiceWg.Add(1)
	go server.Start()
	mLog.Info("(local) internal server > ", interServerIp, ":", interServerPort, " start success.")

	signal.Notify(global.ServiceCh, syscall.SIGINT, syscall.SIGTERM)
	server.Stop()
}
