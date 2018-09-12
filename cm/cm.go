package cm

import (
	"bytes"
	"sort"
	"strconv"
	"strings"

	"github.com/JefferyWang/qcloud/capi"
	"github.com/JefferyWang/qcloud/util"
)

// 云监控请求的地址
const monitorBaseURL = "https://monitor.api.qcloud.com/v2/index.php"

// CM 云监控
type CM struct {
	capi.Conf // 腾讯云api配置
}

// New 实例化监控
func New(secretID string, secretKey string) *CM {
	return &CM{
		Conf: capi.Conf{
			SecretID:  secretID,
			SecretKey: secretKey,
		},
	}
}

// 计算签名
func getSign(method string, reqURI string, params map[string]string) string {
	reqURI = strings.Replace(reqURI, "https://", "", -1)
	secretKey := params["SecretKey"]

	// 字段排序
	var sortedKeys []string
	for key := range params {
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
		v := params[k]
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

	// 计算签名
	return util.HmacSha1Base64(srcStr, secretKey)
}

// GetCvmMonitorData 云服务器监控
func (cmObj *CM) GetCvmMonitorData(serverID string, metric string, period int, startTime string, endTime string) (resp []byte, err error) {
	// 公共参数
	commonParams := cmObj.getCommonParams()
	// 私有或各种类型可能不一致的参数
	params := map[string]string{
		"namespace":          "qce/cvm",
		"dimensions.0.name":  "unInstanceId",
		"dimensions.0.value": serverID,
		"period":             strconv.Itoa(period),
		"Region":             "bj",
	}
	for k, v := range commonParams {
		params[k] = v
	}
	if metric != "" {
		params["metricName"] = metric
	}
	if startTime != "" {
		params["startTime"] = startTime
	}
	if endTime != "" {
		params["endTime"] = endTime
	}
	resp, err = GetMonitorData(params)
	return
}

func (cmObj *CM) getCommonParams() map[string]string {
	return map[string]string{
		"Action":    "GetMonitorData",
		"Timestamp": util.GetCurTimeStampStr(),
		"Nonce":     strconv.FormatUint(util.RandSeed.Uint64(), 10),
		"SecretId":  cmObj.SecretID,
		"SecretKey": cmObj.SecretKey,
	}
}
