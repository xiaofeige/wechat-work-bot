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
参考`example`目录下的实现 
```
// create handler which implement the wx callback msg handler
type ExampleHandler struct{
    *BaseHandler
}

// 覆盖默认的事件响应函数
func (h *ExampleHandler)OnMsgRecv(ctx context.Context, msg *CallBackReq) (rsp *CallBackRsp, err error){
    return &CallbackRsp{}, nil
}

```

添加handler到机器人
```
h := &ExampleHandler{}
robot, err := wechat_work_bot.NewWxWorkRobot(conf.RobotConf)
if err != nil{
    fmt.Println("error")
    return
}
robot.AddHandler(h)
robot.Start()
```

## 如何添砖加瓦
希望我们的合作模式是 fork&pull-request
