package errmsg

const (
	SUCCSE = 200
	ERROR  = 500

	// code= 1000... 用户模块的错误
	ERROR_USERNAME_USED    = 1001
	ERROR_PASSWORD_WRONG   = 1002
	ERROR_USER_NOT_EXIST   = 1003
	ERROR_TOKEN_EXIST      = 1004
	ERROR_TOKEN_RUNTIME    = 1005
	ERROR_TOKEN_WRONG      = 1006
	ERROR_TOKEN_TYPE_WRONG = 1007
	ERROR_USER_NO_RIGHT    = 1008

    //设备相关信息
	ERROR_DEVICE_NOT_EXIST = 4004
	ERROR_DEVICE_EXIST     = 2001
	ERROR_ALL_DEVICE_EXIST = 2002
)

var codeMsg = map[int]string{
	SUCCSE:                 "OK",
	ERROR:                  "FAIL",
	ERROR_USERNAME_USED:    "用户名已存在！",
	ERROR_PASSWORD_WRONG:   "密码错误",
	ERROR_USER_NOT_EXIST:   "用户不存在",
	ERROR_TOKEN_EXIST:      "TOKEN不存在",
	ERROR_TOKEN_RUNTIME:    "TOKEN已过期",
	ERROR_TOKEN_WRONG:      "TOKEN不正确",
	ERROR_TOKEN_TYPE_WRONG: "TOKEN格式错误",
	ERROR_USER_NO_RIGHT:    "该用户无权限",

	ERROR_DEVICE_NOT_EXIST: "设备不存在",
	ERROR_DEVICE_EXIST: "设备已存在",
	ERROR_ALL_DEVICE_EXIST: "所有设备已存在",

}

func GetErrMsg(code int) string {
	return codeMsg[code]
}
