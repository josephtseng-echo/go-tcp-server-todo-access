// global.go ---
//
// Filename: global.go
// Author: josephzeng
// Created:  2018-07-31 19:21:23 (CST)
// Last-Update          2018-08-06 15:55:43 (CST)
//           By:
//     Update #: 0
// Description:全部变量
// Status:
package global

import (
	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
//	"common/server/tcp"
//	"net"
)

var Config = struct {
	ConfigName string
	LogDir     string
	Configer   *ini.File
	Logger     *logrus.Logger
}{}

var ServiceWg sync.WaitGroup = sync.WaitGroup{}
var ServiceCh = make(chan os.Signal, 1)
var ServiceTaskId uint64