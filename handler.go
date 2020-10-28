package wechat_work_bot

import (
	"context"
)

type MsgHandler interface {

	// 机器人id
	RobotId() string

	// 普通文本 & 图片消息处理
	OnMsgRecv(ctx context.Context, msg *CallBackReq) (rsp *CallBackRsp, err error)

	// 事件处理
	// 打开私聊对话框
	OnEnterChat(ctx context.Context, msg *CallBackReq) (rsp *CallBackRsp, err error)

	// 被添加到群聊
	OnAddToChat(ctx context.Context, msg *CallBackReq) (rsp *CallBackRsp, err error)

	// 被踢出群
	OnDeletedFromChat(ctx context.Context, msg *CallBackReq) (rsp *CallBackRsp, err error)

	// 有人加入群聊
	OnGroupAddMember(ctx context.Context, msg *CallBackReq) (rsp *CallBackRsp, err error)

	// 附加事件
	OnAttachmentEvent(ctx context.Context, msg *CallBackReq) (rsp *CallBackRsp, err error)
}

type BaseHandler struct {
}

func (h *BaseHandler) RobotId() string {
	return "default"
}

func (h *BaseHandler) OnMsgRecv(ctx context.Context, msg *CallBackReq) (rsp *CallBackRsp, err error) {
	rsp = &CallBackRsp{
		StrMsgType: WxTextMsg,
		//Markdown:      MarkdownElem{},
		Text: TextRspElem{
			Content: CDATA{Value: "I'm rich!!!"},
			//MentionedList: MentionedListElem{},
		},
		BIgnore: false,
	}

	return rsp, nil
}

func (h *BaseHandler) OnEnterChat(ctx context.Context, msg *CallBackReq) (rsp *CallBackRsp, err error) {

	return &CallBackRsp{BIgnore: true}, nil
}

func (h *BaseHandler) OnAddToChat(ctx context.Context, msg *CallBackReq) (rsp *CallBackRsp, err error) {

	return &CallBackRsp{BIgnore: true}, nil
}

func (h *BaseHandler) OnDeletedFromChat(ctx context.Context, msg *CallBackReq) (rsp *CallBackRsp, err error) {

	return &CallBackRsp{BIgnore: true}, nil
}

func (h *BaseHandler) OnAttachmentEvent(ctx context.Context, msg *CallBackReq) (rsp *CallBackRsp, err error) {

	return &CallBackRsp{BIgnore: true}, nil
}
func (h *BaseHandler) OnGroupAddMember(ctx context.Context, msg *CallBackReq) (rsp *CallBackRsp, err error) {

	return &CallBackRsp{BIgnore: true}, nil
}

//========================================
// 默认的处理
//========================================
type DefaultHandler struct {
	*BaseHandler
}

func (h *DefaultHandler) RobotId() string {
	return "default"
}

func (h *DefaultHandler) OnMsgRecv(ctx context.Context, msg *CallBackReq) (rsp *CallBackRsp, err error) {
	rsp = &CallBackRsp{
		StrMsgType: WxTextMsg,
		//Markdown:      MarkdownElem{},
		Text: TextRspElem{
			Content: CDATA{Value: "I'm rich!!!"},
			//MentionedList: MentionedListElem{},
		},
		BIgnore: false,
	}

	return rsp, nil
}
