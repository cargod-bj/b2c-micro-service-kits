package i18nKits

import (
	"context"
	"fmt"
	"strconv"
)

const defaultSuccess = "0000"

var rcMap = make(map[string]message)

var isDebug = false

// 消息体，对应不同语言
type message struct {
	langEn string
	langId string
	langTh string
	langZh string
}

// 初始化所有微服务的ResponseCode，如果要查询关联lang内容则需要先调用本类方法进行初始化。
// 如果是在微服务中使用，只需要调用相关微服务的Init方法即可，本方法是全量初始化方法。
// 当前是否为debug状态，如果debug=true部分检查可能会报错
func InitRC(allConfig map[string]interface{}, debug bool) {
	isDebug = debug
	println(fmt.Sprintf("initRc:%t", debug))

	initItem := func(target *string, tag string, src map[string]interface{}) {
		if t, find := src[tag]; find {
			switch t.(type) {
			case string:
				*target = t.(string)
			case int:
				i := t.(int)
				*target = strconv.Itoa(i)
			case float64:
				i := t.(float64)
				*target = fmt.Sprintf("%f", i)
			}
		}
	}

	for k := range allConfig {
		msg := message{}
		if m, ok := allConfig[k].(map[string]interface{}); ok {
			initItem(&msg.langEn, "en", m)
			initItem(&msg.langZh, "zh", m)
			initItem(&msg.langTh, "th", m)
			initItem(&msg.langId, "id", m)
		}
		rcMap[k] = msg
	}
}

// 是否成功
func IsSuccess(code string) bool {
	return code == defaultSuccess
}

// 是否失败
func IsFailed(code string) bool {
	return !IsSuccess(code)
}

// 根据responseCode获取responseMessage，使用当前默认语言
// 使用本方法前必须先调用相应微服务的Init类方法
func GetRM(c context.Context, code string) string {
	if isDebug && c == nil {
		panic("the param 'context' is nil")
	}
	return GetRMByLang(code, GetLang(c))
}

// 根据responseCode获取responseMessage，使用lang指定语言
// 使用本方法前必须先调用相应微服务的Init类方法
func GetRMByLang(code, lang string) string {
	if code == "" {
		if isDebug {
			panic("当前传入code为空")
		}
		return "Failed unknown."
	}
	item, ok := rcMap[code]
	if !ok {
		if isDebug {
			panic(fmt.Sprintf("code(%s)不在注册code中，请检查是否错漏", code))
		}
		return "Failed unknown"
	}

	switch lang {
	case LangEn:
		return item.langEn
	case LangId:
		return item.langId
	case LangTh:
		return item.langTh
	case LangZh:
		return item.langZh
	}

	return item.langEn
}
