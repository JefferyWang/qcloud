package cos

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// PutObject 上传本地文件
// 将本地的文件（Object）上传至指定 Bucket 中。该操作需要请求者对 Bucket 有 WRITE 权限。
// objectName: 操作的对象名称
func (cosObj *Cos) PutObject(objectName string, filePath string) (fileURL string, err error) {
	domain := cosObj.getDomain()
	url := fmt.Sprintf("http://%v/%v", domain, objectName)
	headers := map[string]string{
		"Host": domain,
	}
	headers["Authorization"] = cosObj.getAuthorization("put", "/"+objectName, headers, map[string]string{})

	fileByte, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("PUT", url, bytes.NewReader(fileByte))
	if err != nil {
		return
	}
	// 设置头部信息
	for k, v := range headers {
		if k == "Host" {
			req.Host = v
		}
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	// 这个请求如果内容是空，说明是成功的
	if len(respData) == 0 {
		if cosObj.IsHTTPS {
			fileURL = fmt.Sprintf("https://%v/%v", cosObj.FileDomain, objectName)
		} else {
			fileURL = fmt.Sprintf("http://%v/%v", cosObj.FileDomain, objectName)
		}
		return
	}
	err = errors.New(string(respData))
	return
}
