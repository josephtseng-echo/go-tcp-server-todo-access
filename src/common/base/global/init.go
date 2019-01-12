// init.go ---
//
// Filename: init.go
// Author: josephzeng
// Created:  2018-07-31 19:21:58 (CST)
// Last-Update                                                                                              2018-08-06 17:14:07 (CST)
//           By:
//     Update #: 0
// Description: 全局初始化
// Status:

package global

import (
	"common/base/config"
	"common/base/log"
	"os"
)

func InitLogger(name string) {
	loggerCfg := &log.LoggerCfg{
		Dir:  Config.LogDir,
		Name: name,
	}
	Config.Logger = log.NewLogger(loggerCfg)
}

func InitConfiger() (bool, error) {
	_, err := os.Stat(Config.ConfigName)
	if err != nil && os.IsNotExist(err) {
		return false, err
	}
	cfg := config.Configer(Config.ConfigName, "ini")
	if cfg == nil {
		return false, err
	}
	cfg.Section("logs").Key("dir").SetValue(Config.LogDir)
	cfg.SaveTo(Config.ConfigName)
	Config.Configer = cfg
	return true, nil
}
