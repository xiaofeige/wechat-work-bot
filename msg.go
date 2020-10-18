package wechat_work_bot

import (
	"encoding/xml"
)


const (
	WxTextMsg = "text"
	WxMarkdownMsg = "markdown"

	WxEventMsg    = "event"
	WxMixedMsg    = "mixed"
	WxImgMsg      = "image"
)


const (
	EventTypeAddToChat      = "add_to_chat"
	EventTypeDeleteFromChat = "delete_from_chat"
	EventTypeEnterChat      = "enter_chat"
)

//=======================================================
// 收到的加密的消息
//=======================================================
type EncryptedReq struct {
	Encrypt string `xml:"Encrypt"`
}

//=======================================================
// 解密之后的消息
//=======================================================
type MentionedListElem struct {
	Items []CDATA `xml:"Item"`
}

type TextElem struct {
	Content string `xml:"Content"`
}
type ImageElem struct {
	ImageUrl string `xml:"ImageUrl"`
}

type EventElem struct {
	EventType string `xml:"EventType"`
}

type ActionsElem struct {
	Name  string `xml:"Name"`
	Value string `xml:"Value"`
}

type AttachmentElem struct {
	CallbackId string      `xml:"CallbackId"`
	Actions    ActionsElem `xml:"Actions"`
}

type MsgItemElem struct {
	MsgType string    `xml:"MsgType"`
	Text    TextElem  `xml:"Text"`
	Image   ImageElem `xml:"Image"`
}

type MixedMessageElem struct {
	MsgItem MsgItemElem `xml:"MsgItem"`
}

type FromElement struct {
	UserId string `xml:"UserId"`
	Name   string `xml:"Name"`
	Alias  string `xml:"Alias"`
}

type CallBackReq struct {
	MsgItemElem
	WebHookUrl     string         `xml:"WebhookUrl"`
	ChatId         string         `xml:"ChatId"`
	ChatType       string         `xml:"ChatType"`
	GetChatInfoUrl string         `xml:"GetChatInfoUrl"`
	From           FromElement    `xml:"From"`
	MsgId          string         `xml:"MsgId"`
	AppVersion     string         `xml:"AppVersion"`
	Attachment     AttachmentElem `xml:"Attachment"`
	Event          EventElem      `xml:"Event"`
}

func NewWxWorkMessage(data []byte) (*CallBackReq, error) {
	var msg CallBackReq
	err := xml.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

//=======================================================
// 企业微信消息回包
//=======================================================
type CDATA struct {
	Value string `xml:",cdata"`
}
type MarkdownElem struct {
}

type TextRspElem struct {
	Content       CDATA             `xml:"Content"`
	MentionedList MentionedListElem `xml:"MentionedList"`
}
type CallBackRsp struct {
	XMLName       xml.Name     `xml:"xml"`
	StrMsgType    string       `xml:"MsgType"`
	VisibleToUser string       `xml:"VisibleToUser"` // fayren|sumiresun  中线隔开
	Markdown      MarkdownElem `xml:"Markdown"`
	Text          TextRspElem  `xml:"Text"`

	BIgnore bool `xml:"-"`
}

type EncryptedRsp struct {
	XMLName    xml.Name `xml:"xml"`
	EncryptMsg CDATA    `xml:"Encrypt"`
	MsgSign    CDATA    `xml:"MsgSignature"`
	Timestamp  int64    `xml:"TimeStamp"`
	Nonce      CDATA    `xml:"Nonce"`
}
