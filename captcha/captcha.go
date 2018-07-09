package captcha

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/JefferyWang/qcloud/util"
)

// verifyURL 验证码验证地址
const verifyURL = "https://ssl.captcha.qq.com/ticket/verify"

// Conf 验证码配置信息
type Conf struct {
	AppID     string
	SecretKey string
}

// New 实例化
func New(appid string, secretKey string) *Conf {
	return &Conf{
		AppID:     appid,
		SecretKey: secretKey,
	}
}

// respData 验证返回的数据结构
type respData struct {
	Response  string `json:"response"`
	EvilLevel string `json:"evil_level"`
	ErrMsg    string `json:"err_msg"`
}

// Verify 验证码验证
func (conf *Conf) Verify(ticket string, randStr string, userIP string) (ret bool, evilLevel int, err error) {
	resp, err := util.DoGet(verifyURL, map[string]string{
		"aid":          conf.AppID,
		"AppSecretKey": conf.SecretKey,
		"Ticket":       ticket,
		"Randstr":      randStr,
		"UserIP":       userIP,
	})
	if err != nil {
		return
	}

	var verifyData respData
	err = json.Unmarshal(resp, &verifyData)
	if err != nil {
		return
	}

	if verifyData.Response != "1" {
		err = errors.New(verifyData.ErrMsg)
		return
	}

	ret = true
	evilLevel, _ = strconv.Atoi(verifyData.EvilLevel)

	return
}
