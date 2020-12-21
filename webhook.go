package wechat_work_bot

import (
	"fmt"
	"net/url"
)

type RobotWebHook string

// stHookInfo, err := RobotWebHook("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=10ac4f7b-883d-487a-a90d-c504eea860aa").Parse()
func (r RobotWebHook) Parse() (*RobotWebHookInfo, error) {
	var stWebHook RobotWebHookInfo
	err := stWebHook.ParseFrom(string(r))
	if err != nil {
		DefaultLogger.Printf("invalid web hook:%s", string(r))
		return nil, err
	}
	return &stWebHook, nil
}

// strRobotId := RobotWebHook("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=10ac4f7b-883d-487a-a90d-c504eea860aa").RobotId()
func (r RobotWebHook) RobotId() string {

	var stWebHook RobotWebHookInfo
	err := stWebHook.ParseFrom(string(r))
	if err != nil {
		DefaultLogger.Printf("invalid web hook:%s", string(r))
		return ""
	}

	return stWebHook.RobotId
}

type RobotWebHookInfo struct {
	RobotId string `default:""`
	BaseUrl string
}

func (r *RobotWebHookInfo) ParseFrom(hookUrl string) error {

	urlParams, err := url.Parse(hookUrl)
	if err != nil {
		DefaultLogger.Printf("parse webhook url err:%v", err)
		return fmt.Errorf("%w url:%s", ErrInvalidWebHookUrl, hookUrl)
	}
	vals, err := url.ParseQuery(urlParams.RawQuery)
	if err != nil {
		DefaultLogger.Printf("parse raw query err:%v", err)
		return fmt.Errorf("parse raw query err:%w url:%s", ErrInvalidWebHookUrl, hookUrl)
	}
	robotId := vals.Get("key")
	if robotId == "" {
		return fmt.Errorf("%w url:%s", ErrRobotIdNotFound, hookUrl)
	}

	r.RobotId = robotId
	return nil
}
