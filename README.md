# wechat-work-bot
企业微信群聊机器人库

## 功能特性
无需再去看微信文档实现对应的接口，加解密也都帮你实现了。import之后，只需要实现一个handler处理各种消息，然后注册到机器人中即可。
支持注册多个机器人，通过robot id区分。也就是申请的机器人回调中的key的值。

## How to use
```
go get github.com/xiaofeige/wechat-work-bot
```

## example
look into the example dir
```
// create handler which implement the wx callback msg handler
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
```