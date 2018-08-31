package cm

// BaseResp 基础相应数据
type BaseResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
