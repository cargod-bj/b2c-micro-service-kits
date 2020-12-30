package i18nKits

import "fmt"

var localInner string

// 初始化当前全局local
func InitLocal(local string) {
	println("i18n -> current local is", local)
	if !allSupportLocal.Contain(local) {
		panic(fmt.Sprintf("不支持的区域类型：%s", local))
	}
	localInner = local
}

// 获取当前Local信息
func GetLocal() string {
	if localInner == "" {
		panic("请先调用InitLocal方法初始化local信息")
	}
	return localInner
}

// 获取当前local下的默认语言
func getDefaultLangByLocal() string {
	switch localInner {
	case LocalMy:
		return LangEn
	case LocalId:
		return LangId
	case LocalIn:
		return LangEn
	case LocalTh:
		return LangTh
	case LocalZh:
		return LangZh
	default:
		return LangEn
	}
}
