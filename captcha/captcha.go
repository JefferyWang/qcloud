package captcha

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

// Verify 验证码验证
func Verify(ticket string, randStr string, userIP string) (ret bool, evilLevel int, err error) {
	// TODO
	return
}
