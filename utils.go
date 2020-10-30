package wechat_work_bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
)

func StructToJson(obj interface{}) string {
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.SetIndent("", "	")
	jsonEncoder.Encode(obj)
	return bf.String()
}

func GetRobotIdByHookUrl(hookUrl string) (string, error) {
	urlParams, err := url.Parse(hookUrl)
	if err != nil {
		DefaultLogger.Printf("parse webhook url err:%v", err)
		return "", fmt.Errorf("%w url:%s", ErrInvalidWebHookUrl, hookUrl)
	}
	vals, err := url.ParseQuery(urlParams.RawQuery)
	if err != nil {
		DefaultLogger.Printf("parse raw query err:%v", err)
		return "", fmt.Errorf("parse raw query err:%w url:%s", ErrInvalidWebHookUrl, hookUrl)
	}
	robotId := vals.Get("key")
	if robotId == "" {
		return "", fmt.Errorf("%w url:%s", ErrRobotIdNotFound, hookUrl)
	}
	return robotId, nil
}
