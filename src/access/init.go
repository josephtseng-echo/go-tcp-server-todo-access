package main

import (
	"common/base/global"
	"fmt"
	"github.com/urfave/cli"
	"os"
)

func InitServer() {
	c, err := global.InitConfiger()
	if c == false {
		fmt.Println("配置有问题 ERROR:" + err.Error())
		os.Exit(1)
	}
	global.InitLogger("logserver")
}

func initArgs() {
	app := cli.NewApp()
	app.Version = "0.0.1"
	app.Name = "ACCESS服务"
	app.Usage = "ACCESS网关服务"
	app.UsageText = ""
	app.ArgsUsage = "参数的描述"

	app.Email = "josephzeng@local"
	app.Author = "josephzeng"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "配置文件绝对路径，如： /data/config.ini; 必须是ini文件",
		},
		cli.StringFlag{
			Name:  "logs, l",
			Usage: "配置logs根路径，绝对路径，如： /data/logs ",
			Value: "/data/logs",
		},
	}

	app.Action = func(c *cli.Context) error {
		global.Config.LogDir = c.String("logs")
		if c.String("config") != "" {
			global.Config.ConfigName = c.String("config")
			InitServer()
			return nil
		}
		cli.ShowAppHelp(c)
		os.Exit(1)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
	}
}
