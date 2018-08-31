package cm

import (
	"github.com/JefferyWang/qcloud/util"
)

// GetMonitorData 获取监控数据
func GetMonitorData(params map[string]string) (resp []byte, err error) {
	// 获取签名
	params["Signature"] = getSign("GET", monitorBaseURL, params)
	// 删除secretKey，避免泄露
	delete(params, "SecretKey")

	headers := map[string]string{}
	resp, err = util.DoGet(monitorBaseURL, params, headers)
	if err != nil {
		return
	}
	return
}
