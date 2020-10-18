package wechat_work_bot

import (
	"bytes"
	"encoding/json"
)


func StructToJson(obj interface{}) string {
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.SetIndent("", "	")
	jsonEncoder.Encode(obj)
	return bf.String()
}
