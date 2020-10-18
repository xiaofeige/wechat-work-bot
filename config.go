package wechat_work_bot



type RobotServerConfig struct {
	Bind 			string 	`default:"127.0.0.1"`
	Port 			int		`default:"7788"`
	UrlBase 		string 	`default:"/"`
}