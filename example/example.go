package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	wechat_work_bot "github.com/wechat-work-bot"
)

type ExampleConfig struct {
	Business		struct{
		// business config
	}

	RobotConf 		wechat_work_bot.RobotConfig
}

func main(){

	h, err := NewExamplerHandler()
	if err != nil{
		fmt.Println("new example handler err:", err)
		return
	}

	var conf ExampleConfig
	_, err = toml.Decode("./config.toml", &conf)
	if err != nil{
		fmt.Println("parse config err:", err)
		return
	}

	robot, err := wechat_work_bot.NewWxWorkRobot(conf.RobotConf)
	if err != nil{
		fmt.Println("")
	}
	robot.AddHandler(h)
	fmt.Println("wechat-work-bot exit with:", err)
}
