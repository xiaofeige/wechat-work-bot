package wechat_work_bot

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type RobotConfig struct {
	ReceiveId string `default:""` // 在群聊机器人场景中，直接传空字符串

	// 申请机器人时的 token
	Token string `default:""`
	// 申请机器人的时候的base64 编码的aeskey
	TokenAESKey string
	// base64 解码之后的aeskey, 传入config的时候不需要指定
	aesKey []byte

	WebHookUrl string //企业微信推送消息地址

	BindAddr string `default:":8080"` //

	UrlBase string `default:"/"`

	// 消息处理接口
	Handler MsgHandler

	//
	Debugger Logger

	//
	ErrLogger Logger
}

type Robot struct {
	conf *RobotConfig
	*gin.Engine

	handlers map[string]MsgHandler

	//
	Debugger Logger

	//
	ErrLogger Logger
}

// 创建wx-work robot回调机器人
func NewWxWorkRobot(conf RobotConfig) (*Robot, error) {

	aseKey, err := base64.StdEncoding.DecodeString(conf.TokenAESKey + "=") // 不知道为啥这样设计，要这里自己手动添加等号
	if err != nil {
		fmt.Printf("aes key[%s] base64 decode err:%v\n", conf.TokenAESKey, err)
		return nil, err
	}
	conf.aesKey = aseKey

	if conf.Handler == nil {
		conf.Handler = &DefaultHandler{}
	}

	if conf.Debugger == nil {
		conf.Debugger = DefaultLogger
	}

	if conf.ErrLogger == nil {
		conf.ErrLogger = DefaultLogger
	}

	robot := &Robot{
		conf:   &conf,
		Engine: gin.Default(),
	}

	return robot, nil
}

func (r *Robot) AddHandler(handler MsgHandler) {
	r.handlers[handler.RobotId()] = handler
}

func (r *Robot) Start() {
	r.Handle(http.MethodGet, r.conf.UrlBase, r.VerifyUrl)
	r.Handle(http.MethodPost, r.conf.UrlBase, r.OnRecvMsg)

	r.Debugger.Printf("running config:\n %s", StructToJson(r.conf))
	err := r.Run(r.conf.BindAddr)
	if err != nil {
		r.ErrLogger.Printf("robot exit...err:%v", err)
	}
}

func (r *Robot) Encrypt(origData []byte) (string, error) {

	var buf []byte
	// 加入随机字符串
	buf = append(buf, []byte(getRandomStr(16))...)
	// 标记数据长度
	lenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(origData)))
	buf = append(buf, lenBytes...)
	// 压入数据
	buf = append(buf, origData...)
	// 填充received id
	if len(r.conf.ReceiveId) > 0 {
		buf = append(buf, []byte(r.conf.ReceiveId)...)
	}
	return Encrypt(buf, r.conf.aesKey)
}

func (r *Robot) Decrypt(encryptedData string) ([]byte, error) {

	plainText, err := Decrypt(encryptedData, r.conf.aesKey)
	if err != nil {
		return nil, err
	}

	msgLenBytes := plainText[16:20]

	msgLen := binary.BigEndian.Uint32(msgLenBytes)

	if int(msgLen)+20 > len(plainText) {
		return nil, fmt.Errorf("invalid msg len:%d", msgLen)
	}
	msg := plainText[20 : msgLen+20]

	return msg, nil
}

func (r *Robot) ValidSign(strSign, strTs, strNonce, strData string) error {
	iTs, err1 := strconv.ParseInt(strTs, 10, 64)
	if err1 != nil {
		return fmt.Errorf("err1:%v", err1)
	}

	localSign := Sha1(r.conf.Token, strData, strNonce, iTs)
	if strSign != localSign {
		return fmt.Errorf("invalid sign")
	}
	return nil
}

func (r *Robot) VerifyUrl(ctx *gin.Context) {
	sign := ctx.Query("msg_signature")
	ts := ctx.Query("timestamp")
	nonce := ctx.Query("nonce")
	strEcho := ctx.Query("echostr")

	if err := r.ValidSign(sign, ts, nonce, strEcho); err != nil {
		return
	}

	msg, err := r.Decrypt(strEcho)
	if err != nil {
		return
	}

	ctx.String(http.StatusOK, string(msg))
}

func (r *Robot) OnRecvMsg(ctx *gin.Context) {

	// 接收数据 & 校验数据
	sign := ctx.Query("msg_signature")
	ts := ctx.Query("timestamp")
	nonce := ctx.Query("nonce")

	var req EncryptedReq
	err := ctx.BindXML(&req)
	if err != nil {
		r.ErrLogger.Printf("parse request err:%v", err)
		return
	}

	if err := r.ValidSign(sign, ts, nonce, req.Encrypt); err != nil {
		r.ErrLogger.Printf("sign check err:%v", err)
		return
	}

	msgData, err := r.Decrypt(req.Encrypt)
	if err != nil {
		r.ErrLogger.Printf("parse encrypted msg err:%v", err)
		return
	}
	r.Debugger.Printf("msgData:%s", string(msgData))

	msg, err := NewWxWorkMessage(msgData)
	if err != nil {
		r.ErrLogger.Printf("parse new msg err:%v", err)
		return
	}
	r.Debugger.Printf("recv msg:%s", StructToJson(msg))

	// 分发事件
	rsp, err := r.HandleMsg(ctx, msg)
	if err != nil {
		r.ErrLogger.Printf("handle msg err:%v", err)
		return
	}
	if rsp.BIgnore {
		r.Debugger.Printf("ignore rsp...")
		return
	}

	rspBytes, err := xml.Marshal(rsp)
	if err != nil {
		r.ErrLogger.Printf("marshal rsp err:%v", err)
		return
	}
	r.Debugger.Printf("rsp xml:%s", string(rspBytes))

	strEncryptedRsp, err := r.Encrypt(rspBytes)
	if err != nil {
		r.ErrLogger.Printf("encrypt rsp err:%v", err)
		return
	}
	nonce = fmt.Sprint(time.Now().Unix())
	iTs := time.Now().Unix()
	sign = Sha1(r.conf.Token, strEncryptedRsp, nonce, iTs)

	encryptedRsp := &EncryptedRsp{
		EncryptMsg: CDATA{Value: strEncryptedRsp},
		MsgSign:    CDATA{sign},
		Timestamp:  iTs,
		Nonce:      CDATA{nonce},
	}

	strRsp, err := xml.Marshal(encryptedRsp)
	if err != nil {
		r.ErrLogger.Printf("marshal encrypted rsp err:%v", err)
		return
	}
	r.Debugger.Printf("str rsp:%s", strRsp)

	ctx.String(http.StatusOK, string(strRsp))
}

func (r *Robot) HandleMsg(ctx context.Context, req *CallBackReq) (rsp *CallBackRsp, err error) {
	robotId, err := GetRobotIdByHookUrl(req.WebHookUrl)
	if err != nil {
		r.ErrLogger.Printf("get robot id from callback req err:%v", err)
		return nil, err
	}

	handler, exist := r.handlers[robotId]
	if !exist {
		r.ErrLogger.Printf("robot id:%s handler is not found, using default handler", robotId)
		handler = &DefaultHandler{}
	}
	switch req.MsgType {
	case WxTextMsg:
		r.Debugger.Printf("recv text msg...")
		return handler.OnMsgRecv(ctx, req)
	case WxEventMsg:
		switch req.Event.EventType {
		case EventTypeAddToChat:
			return handler.OnAddToChat(ctx, req)
		case EventTypeDeleteFromChat:
			return handler.OnDeletedFromChat(ctx, req)
		case EventTypeEnterChat:
			return handler.OnEnterChat(ctx, req)
		default:
			r.ErrLogger.Printf("unsupported event type:%s", req.Event.EventType)
			return &CallBackRsp{BIgnore: true}, nil
		}
	default:
		r.ErrLogger.Printf("unsupported msg type:%s", req.MsgType)
		return &CallBackRsp{BIgnore: true}, nil
	}
}
