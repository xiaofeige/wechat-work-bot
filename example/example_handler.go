package main

import (
	"context"
	wechat_work_bot "github.com/xiaofeige/wechat-work-bot"
)

type ExampleHandler struct {
	*wechat_work_bot.BaseHandler
}

func NewExampleHandler() (*ExampleHandler, error) {
	h := &ExampleHandler{}
	return h, nil
}

func (h *ExampleHandler) RobotId() string {
	// webhook: https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=da1b0595-cd46-4b43-bc3a-b04e0e3a2bca
	return "da1b0595-cd46-4b43-bc3a-b04e0e3a2bca"
}

func (h *ExampleHandler) OnMsgRecv(ctx context.Context, msg *wechat_work_bot.CallBackReq) (rsp *wechat_work_bot.CallBackRsp, err error) {
	rsp = &wechat_work_bot.CallBackRsp{
		StrMsgType: wechat_work_bot.WxTextMsg,
		//Markdown:      MarkdownElem{},
		Text: wechat_work_bot.TextRspElem{
			Content: wechat_work_bot.CDATA{Value: "I'm rich!!!"},
			//MentionedList: MentionedListElem{},
		},
		BIgnore: false,
	}

	return rsp, nil
}

func (h *ExampleHandler) OnEnterChat(ctx context.Context, msg *wechat_work_bot.CallBackReq) (rsp *wechat_work_bot.CallBackRsp, err error) {

	return &wechat_work_bot.CallBackRsp{BIgnore: true}, nil
}

func (h *ExampleHandler) OnAddToChat(ctx context.Context, msg *wechat_work_bot.CallBackReq) (rsp *wechat_work_bot.CallBackRsp, err error) {

	return &wechat_work_bot.CallBackRsp{BIgnore: true}, nil
}

func (h *ExampleHandler) OnDeletedFromChat(ctx context.Context, msg *wechat_work_bot.CallBackReq) (rsp *wechat_work_bot.CallBackRsp, err error) {

	return &wechat_work_bot.CallBackRsp{BIgnore: true}, nil
}

func (h *ExampleHandler) OnAttachmentEvent(ctx context.Context, msg *wechat_work_bot.CallBackReq) (rsp *wechat_work_bot.CallBackRsp, err error) {

	return &wechat_work_bot.CallBackRsp{BIgnore: true}, nil
}
