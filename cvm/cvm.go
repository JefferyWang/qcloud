package cvm

import (
	"bytes"
	"sort"
	"strconv"
	"strings"

	"github.com/JefferyWang/qcloud/capi"
	"github.com/JefferyWang/qcloud/util"
)

// REQUEST_URL 请求地址，该域名支持就近接入，不需要指定地域
const REQUEST_URL = "https://cvm.tencentcloudapi.com/"

// CVM 云服务器
type CVM struct {
	capi.Conf
	Region string
}

// Params 请求参数
type Params map[string]string

// New 示例CVM接口
func New(secretID string, secretKey string, region string) *CVM {
	return &CVM{
		Conf: capi.Conf{
			SecretID:  secretID,
			SecretKey: secretKey,
		},
		Region: region,
	}
}

// NewParams 生成公共参数
func (cvmObj *CVM) NewParams() *Params {
	return &Params{
		"Timestamp":       util.GetCurTimeStampStr(),
		"Nonce":           strconv.FormatUint(util.RandSeed.Uint64(), 10),
		"SecretId":        cvmObj.SecretID,
		"SecretKey":       cvmObj.SecretKey,
		"Region":          cvmObj.Region,
		"Version":         "2017-03-12",
		"SignatureMethod": "HmacSHA1",
	}
}

// GetSign 获取签名
func (params *Params) GetSign(method string, reqURI string) {
	reqURI = strings.Replace(reqURI, "https://", "", -1)

	secretKey := (*params)["SecretKey"]

	// 字段排序
	var sortedKeys []string
	for key := range *params {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)

	// 字段拼接
	var buf bytes.Buffer
	for _, k := range sortedKeys {
		// 忽略部分字段
		if k == "SecretKey" || k == "Signature" {
			continue
		}
		v := (*params)[k]
		if v != "" {
			if buf.Len() > 0 {
				buf.WriteString(`&`)
			}
			buf.WriteString(k)
			buf.WriteString(`=`)
			buf.WriteString(v)
		}
	}
	srcStr := method + reqURI + "?" + buf.String()

	(*params)["Signature"] = util.HmacSha1Base64(srcStr, secretKey)
	delete(*params, "SecretKey")
}
