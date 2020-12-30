package i18nKits

import (
	"context"
	"fmt"
	"github.com/cargod-bj/b2c-common/csErr"
)

// 新建一个rc类型的错误
// args 用来format code对应的message，例如code对应的message是：%s is error, 入参规则同fmt的format
func Err(c context.Context, code string, args ...interface{}) csErr.CSError {
	var m = GetRM(c, code)
	if len(args) > 0 {
		m = fmt.Sprintf(m, args...)
	}

	e := csErr.CreateError(nil, code, m, GetRMByLang(code, LangEn))

	return e
}

// 新建一个rc类型的错误, 使用给定msg作为rm
func ErrWithMsg(c context.Context, code string, msg string) csErr.CSError {
	var m = GetRM(c, code)
	if msg != "" {
		m = msg
	}

	e := csErr.CreateError(nil, code, m, GetRMByLang(code, LangEn))

	return e
}

// 新建一个rc类型的错误, 使用给定msg作为rm
func ErrWithParent(c context.Context, code string, parent error) csErr.CSError {
	var m = GetRM(c, code)

	e := csErr.CreateError(parent, code, m, GetRMByLang(code, LangEn))

	return e
}

//
//// 将给定error转换为RCError类型。
//// 如果转换后的error中的rc或rm为空，则使用codeMsg参数进行填充
//// codeMsg 接收两个参数，第一个参数作为code使用，第二个参数作为message使用
//func errTo(c *context.Context, e error, codeMsg ...string) *rcError {
//	createError := func() *rcError {
//		var rc = ""
//		var rm = ""
//		if len(codeMsg) > 0 {
//			rc = codeMsg[0]
//		}
//		if len(codeMsg) > 1 {
//			rm = codeMsg[1]
//		}
//		rce := rcError{rc: rc, rm: rm, originalMsg: e.Error()}
//		return &rce
//	}
//
//	if e == nil {
//		return createError()
//	}
//
//	switch e.(type) {
//	case *rcError:
//		er := e.(*rcError)
//		if er.rc == "" && len(codeMsg) > 0 {
//			er.rc = codeMsg[0]
//		}
//		if er.rm == "" && len(codeMsg) > 1 {
//			er.rm = codeMsg[1]
//		}
//		return er
//	default:
//		return createError()
//	}
//}
//
//// 将给定error转换为RCError类型
//// 如果传入了codeMsg参数，则替换原error中的相应参数
//// codeMsg 接收两个参数，第一个参数作为code使用，第二个参数作为message使用
//func errReplace(c *context.Context, e error, codeMsg ...string) *rcError {
//	createError := func() *rcError {
//		var rc = ""
//		var rm = ""
//		if len(codeMsg) > 0 {
//			rc = codeMsg[0]
//		}
//		if len(codeMsg) > 1 {
//			rm = codeMsg[1]
//		}
//		rce := rcError{rc: rc, rm: rm, originalMsg: e.Error()}
//		return &rce
//	}
//
//	if e == nil {
//		return createError()
//	}
//
//	switch e.(type) {
//	case *rcError:
//		return e.(*rcError)
//	default:
//		return createError()
//	}
//}
//
//// 获取response code
//func (e *rcError) GetRC() string {
//	return e.rc
//}
//
//// 获取response message
//func (e *rcError) GetRM() string {
//	return e.rm
//}
//
//// 返回详细的错误信息，如果有原始错误信息，也会输出原始错误信息
//func (e *rcError) Error() string {
//	om := ""
//	if e.originalMsg != "" {
//		om = " errMsg: " + e.originalMsg
//	}
//	return "(" + e.rc + ")" + e.prm + om
//}
//
//// 返回主要错误信息，不返回原始错误信息
//func (e *rcError) String() string {
//	return e.rm + " (" + e.rc + ")"
//}
