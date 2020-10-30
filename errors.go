package wechat_work_bot

import "errors"

var (
	ErrRobotIdNotFound = errors.New("robot id not found")

	ErrInvalidWebHookUrl = errors.New("invalid web hook url")
)
