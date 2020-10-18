package wechat_work_bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//===========================================
// 本包主要包含企业微信接口api
//===========================================
func SendRobotMsg(msgType string, alarmUrl, content string) error{

	msgData := make(map[string]interface{})

	msgData["msgtype"] = msgType
	msgData[string(msgType)] = map[string]interface{}{
		"content": content,
	}

	jsonData, err := json.Marshal(msgData)
	if err != nil {
		return err
	}
	return sendRobotMsg(alarmUrl, jsonData)
}

func sendRobotMsg(alarmUrl string, jsonData []byte) error {
	req, err := http.NewRequest(http.MethodPost, alarmUrl, bytes.NewReader(jsonData))
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