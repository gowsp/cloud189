package invoker

import "fmt"

// 错误码
const (
	ErrFileNotFound           = "FileNotFound"           //文件不存在
	ErrInvalidArgument        = "InvalidArgument"        //参数错误
	ErrUserDayFlowOverLimited = "UserDayFlowOverLimited" //用户当日流量已用完
)

var invalid = map[string]struct{}{
	"InvalidSignature":   {},
	"InvalidSessionKey":  {},
	"InvalidAccessToken": {},
}

type OkRsp interface {
	error
	IsSuccess() bool
}
type NumCodeRsp struct {
	ResCode    int    `json:"res_code"`
	ResMessage string `json:"res_message"`
}

func (r *NumCodeRsp) IsSuccess() bool {
	return r.ResCode == 0
}
func (r *NumCodeRsp) Error() string {
	return fmt.Sprintf("%d: %s", r.ResCode, r.ResMessage)
}
func (r *NumCodeRsp) Code() int {
	return r.ResCode
}
func (r *NumCodeRsp) Message() string {
	return r.ResMessage
}

type BadRsp interface {
	error
	IsError(string) bool
}

type strCodeRsp struct {
	StrCode    string `json:"code"`
	StrMessage string `json:"msg"`
	ResCode    string `json:"res_code"`
	ResMessage string `json:"res_message"`
}

// 业务异常, 无需重试
func (r *strCodeRsp) isBusinessErr() bool {
	_, ok := invalid[r.Code()]
	return !ok
}
func (r *strCodeRsp) IsError(code string) bool {
	return r.Code() == code
}
func (r *strCodeRsp) Error() string {
	return fmt.Sprintf("%s: %s", r.Code(), r.Message())
}
func (r *strCodeRsp) Code() string {
	if r.ResCode != "" {
		return r.ResCode
	}
	return r.StrCode
}
func (r *strCodeRsp) Message() string {
	if r.ResMessage != "" {
		return r.ResMessage
	}
	return r.StrMessage
}
