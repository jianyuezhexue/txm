package util

import (
	"runtime"
	"strings"
)

// 获取函数名
func GetFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	nm := runtime.FuncForPC(pc).Name()
	return nm[strings.LastIndex(nm, ".")+1:]
}
