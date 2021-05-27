package wechat_work_bot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//===========================================
// 本包主要包含企业微信接口api
//===========================================
func SendRobotMsg(msgType string, strWebHookUrl, content string) error {

	msgData := make(map[string]interface{})

	msgData["msgtype"] = msgType
	msgData[string(msgType)] = map[string]interface{}{
		"content": content,
	}

	jsonData, err := json.Marshal(msgData)
	if err != nil {
		return err
	}
	return sendRobotMsg(strWebHookUrl, jsonData)
}

// 通用发送消息函数, 可以参考下述方式发送消息，简单易实现
//content := map[string]interface{}{
//		"msgtype": "markdown",
//		//"visible_to_user": "fayren", //诶，就是不开启这个
//		"markdown": map[string]interface{}{
//			"content": "",
//			"at_short_name": true,
//			"attachments": []map[string]interface{}{
//				map[string]interface{}{
//					"callback_id": "on_click_confirm",
//					"actions": []map[string]interface{}{
//						map[string]interface{}{
//							"name": "agree",
//							"text": "准奏",
//							"type": "button",
//							"value": "agree",
//							"replace_text": "谢主隆恩",
//							"border_color": "2EAB49",
//							"text_color": "2EAB49",
//						},
//
//						map[string]interface{}{
//							"name": "reject",
//							"text": "此事再议",
//							"type": "button",
//							"value": "reject",
//							"replace_text": "谢主隆恩",
//							"border_color": "2EAB49",
//							"text_color": "2EAB49",
//						},
//					},
//				},
//			},
//		},
//	}
func SendWxRobotMsg(strWebHookUrl string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return sendRobotMsg(strWebHookUrl, jsonData)
}

func sendRobotMsg(strWebHookUrl string, jsonData []byte) error {
	req, err := http.NewRequest(http.MethodPost, strWebHookUrl, bytes.NewReader(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	cli := &http.Client{}
	rsp, err := cli.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	bodyData, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}
	if rsp.StatusCode != http.StatusOK {
		return fmt.Errorf("call send msg err, http code:%d, body:%s", rsp.StatusCode, string(bodyData))
	}
	return nil
}

//获取群聊资料
func GetGroupChatInfo(ctx context.Context, groupInfoUrl string) (*GroupInfoRsp, error) {

	httpRsp, err := http.Get(groupInfoUrl)
	if err != nil {
		return nil, err
	}

	bodyData, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		return nil, err
	}

	var rsp GroupInfoRsp
	err = json.Unmarshal(bodyData, &rsp)
	if err != nil {
		return nil, err
	}

	return &rsp, nil
}
