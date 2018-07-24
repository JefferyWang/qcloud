package cos

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"io"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/JefferyWang/qcloud/capi"
)

// Cos 对象存储
type Cos struct {
	capi.Conf         // 腾讯云api配置
	Name       string // 对象存储的bucket name
	BucketID   string // 对象存储的bucket appid
	Region     string // 对象存储所在的区域，可用区域参考https://cloud.tencent.com/document/product/436/6224
	SignExpire uint   // 签名过期时间，单位秒
	FileDomain string // 最终文件下载的域名，如果开启CDN，并使用CDN，此处请填写CDN域名
	IsHTTPS    bool   // 最终url是否要使用https
}

// New 实例化
func New(appid string, secretID string, secretKey string, bucketName string, bucketID string, region string) *Cos {
	cosObj := &Cos{
		Conf: capi.Conf{
			AppID:     appid,
			SecretID:  secretID,
			SecretKey: secretKey,
		},
		Name:       bucketName,
		BucketID:   bucketID,
		Region:     region,
		SignExpire: 600,
		IsHTTPS:    false,
	}
	cosObj.FileDomain = cosObj.getDomain()
	return cosObj
}

// SetHTTPSFlag 设置是否支持https
func (cosObj *Cos) SetHTTPSFlag(flag bool) *Cos {
	cosObj.IsHTTPS = flag
	return cosObj
}

// getDomain 获取接口域名
func (cosObj *Cos) getDomain() string {
	// 格式： {bucket-appid}.cos.{region}.myqcloud.com
	// 示例： zuhaotestsgnoversion-1251668577.cos.ap-beijing.myqcloud.com
	return fmt.Sprintf("%v-%v.cos.%v.myqcloud.com", cosObj.Name, cosObj.BucketID, cosObj.Region)
}

// getAuthorization 获取请求签名信息
func (cosObj *Cos) getAuthorization(method string, uri string, headers map[string]string, params map[string]string) string {
	headerKeyStr, headerDataStr := getListStr(headers)
	paramKeyStr, paramDataStr := getListStr(params)
	data := map[string]string{
		"method":           method,
		"uri":              uri,
		"q-sign-algorithm": "sha1",                  // 签名算法，目前仅支持 sha1，即为 sha1 。
		"q-ak":             cosObj.SecretID,         // 帐户 ID，即 SecretId
		"q-sign-time":      cosObj.getSignTimeStr(), // 本签名的有效起止时间，以秒为单位，格式为 [start-seconds];[end-seconds]。
		"q-key-time":       "",                      // 与 q-sign-time 值相同。
		"q-header-list":    headerKeyStr,            // HTTP 请求头部。需从 key:value 中提取部分或全部 key，且 key 需转化为小写，并将多个 key 之间以字典顺序排序，如有多组 key，可用“;”连接。
		"q-url-param-list": paramKeyStr,             // HTTP 请求参数。需从 key=value 中提取部分或全部 key，且 key 需转化为小写，并将多个 key 之间以字典顺序排序，如有多组 key，可用“;”连接。
		"q-signature":      "",                      // HTTP内容签名
		"header-data-str":  headerDataStr,
		"param-data-str":   paramDataStr,
	}
	data["q-key-time"] = data["q-sign-time"]
	data["q-signature"] = cosObj.getSign(data)
	requiredFields := []string{
		"q-sign-algorithm", "q-ak", "q-sign-time", "q-key-time", "q-header-list",
		"q-url-param-list", "q-signature",
	}
	var buf bytes.Buffer
	for _, k := range requiredFields {
		v := data[k]
		if buf.Len() > 0 {
			buf.WriteString(`&`)
		}
		buf.WriteString(k)
		buf.WriteString(`=`)
		buf.WriteString(v)
	}
	return buf.String()
}

// getSign 获取请求内容签名
func (cosObj *Cos) getSign(data map[string]string) string {
	// 签名算法
	// 第一步： 对临时密钥的有效起止时间加密计算值 SignKey。
	signKey := hmacSha1(data["q-key-time"], cosObj.SecretKey)
	// 第二步： 根据固定格式组合生成 HttpString。
	// [HttpMethod]\n[HttpURI]\n[HttpParameters]\n[HttpHeaders]\n
	httpString := fmt.Sprintf("%v\n%v\n%v\n%v\n", data["method"], data["uri"], data["param-data-str"], data["header-data-str"])
	sha1edHTTPString := doSha1(httpString)
	// 第三步： 加密 HttpString，并根据固定格式组合生成 StringToSign。
	// [q-sign-algorithm]\n[q-sign-time]\nsha1($HttpString)\n
	stringToSign := fmt.Sprintf("%v\n%v\n%v\n", data["q-sign-algorithm"], data["q-sign-time"], sha1edHTTPString)
	// 第四步： 加密 StringToSign，生成Signature。
	return hmacSha1(stringToSign, signKey)
}

// getSignTimeStr 获取签名有效起止时间字符串
func (cosObj *Cos) getSignTimeStr() string {
	now := time.Now().Unix()
	end := now + int64(cosObj.SignExpire)

	return fmt.Sprintf("%v;%v", now, end)
}

// getListStr 获取列表头部字符串
// 返回两个字段
// 第一个： 是list中key的列表
// 第二个： 是所有信息的连接数据
func getListStr(data map[string]string) (string, string) {
	if len(data) <= 0 {
		return "", ""
	}
	// 先对数据转小写
	var keys = make([]string, len(data))
	var tmpData = make(map[string]string)
	i := 0
	for k, v := range data {
		tmpKey := strings.ToLower(k)
		keys[i] = tmpKey
		tmpData[tmpKey] = url.QueryEscape(v)
		i++
	}
	// 再对数据做排序
	sort.Strings(keys)
	var buf bytes.Buffer
	for _, k := range keys {
		v := tmpData[k]
		if buf.Len() > 0 {
			buf.WriteString(`&`)
		}
		buf.WriteString(k)
		buf.WriteString(`=`)
		buf.WriteString(v)
	}
	return strings.Join(keys, ";"), buf.String()
}

func hmacSha1(data, key string) string {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(data))
	return fmt.Sprintf("%x", mac.Sum(nil))
}

func doSha1(data string) string {
	h := sha1.New()
	io.WriteString(h, data)
	return fmt.Sprintf("%x", h.Sum(nil))
}
