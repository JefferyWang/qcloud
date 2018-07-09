package util

import (
	"io/ioutil"
	"net/http"
)

// DoGet 发送get请求
func DoGet(urlStr string, params map[string]string) (body []byte, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return
	}

	// 添加get请求参数
	query := req.URL.Query()
	for key, val := range params {
		query.Set(key, val)
	}
	req.URL.RawQuery = query.Encode()

	// 发起请求
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return
}
