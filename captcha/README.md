# captcha

## 说明

腾讯验证码(腾讯防水墙) Golang版本SDK

官网地址: [https://open.captcha.qq.com/index.html](https://open.captcha.qq.com/index.html)

官方文档地址： [https://007.qq.com/captcha/#/gettingStart?ADTAG=index.head](https://007.qq.com/captcha/#/gettingStart?ADTAG=index.head)

## 快速使用

``` go
import (
	"fmt"
	"github.com/JefferyWang/qcloud/captcha"
)

func main() {
	ticket := "****"
	randStr := "****"
	userIP := "*.*.*.*"

	ret, level, err := captcha.New("your appid", "your secret key").Verify(ticket, randStr, userIP)
	if err != nil {
		fmt.Println("[Error]: ", err)
		return
	}
	fmt.Println(ret, level, err)
}
```



