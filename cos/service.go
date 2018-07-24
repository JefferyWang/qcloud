package cos

import (
	"encoding/xml"
	"fmt"

	"github.com/JefferyWang/qcloud/util"
)

// GetService 获取请求者名下的所有存储空间列表
func (cosObj *Cos) GetService() (result BucketsResult, err error) {
	urlStr := "http://service.cos.myqcloud.com"
	headers := map[string]string{
		"Host": "service.cos.myqcloud.com",
	}
	headers["Authorization"] = cosObj.getAuthorization("get", "/", headers, map[string]string{})
	params := map[string]string{}
	resp, err := util.DoGet(urlStr, params, headers)
	if err != nil {
		return
	}

	err = xml.Unmarshal(resp, &result)
	if err != nil {
		return
	}
	if result.Code != "" {
		err = fmt.Errorf("Get Service failed, code: %v, msg: %v, requestId: %v, traceId: %v", result.Code, result.Message, result.RequestID, result.TraceID)
		return
	}
	return
}
