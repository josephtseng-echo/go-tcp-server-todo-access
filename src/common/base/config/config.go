// config.go ---
//
// Filename: config.go
// Author: josephzeng
// Created:  2018-07-17 15:52:27 (CST)
// Last-Update                                                                                                                                                                                                                              2018-07-17 16:32:01 (CST)
//           By:
//     Update #: 0
// Description:
// Status:
package config

import (
	"github.com/go-ini/ini"
)

type Config interface {
	loadInit() *ini.File
}

type config struct {
	filetype string
	filename string
}

func (self *config) loadInit() *ini.File {
	cfg, err := ini.Load(self.filename)
	if err != nil {
		return nil
	}
	return cfg
}

func Configer(filename string, filetype string) *ini.File {
	cfg := &config{
		filetype: filetype,
		filename: filename,
	}
	if filetype == "ini" {
		return cfg.loadInit()
	}
	return nil
}
