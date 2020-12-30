package i18nKits

import (
	"context"
	setKits "github.com/cargod-bj/b2c-common/utils/set"
	"strings"
)

const (
	// 英语
	LangEn = "en"
	// 印尼语
	LangId = "id"
	// 泰语
	LangTh = "th"
	// 汉语
	LangZh = "zh"
)

const (
	// 马来
	LocalMy = "my"
	// 印尼
	LocalId = "id"
	// 印度
	LocalIn = "in"
	// 泰国
	LocalTh = "th"
	// 中国
	LocalZh = "zh"
)

var allSupportLocal setKits.ISetStr
var allSupportLang setKits.ISetStr

func init() {
	ss := setKits.SetStr()
	allSupportLocal = &ss
	allSupportLocal.Add(LocalMy)
	allSupportLocal.Add(LocalId)
	allSupportLocal.Add(LocalIn)
	allSupportLocal.Add(LocalTh)
	allSupportLocal.Add(LocalZh)

	ss = setKits.SetStr()
	allSupportLang = &ss

	allSupportLang.Add(LangEn)
	allSupportLang.Add(LangId)
	allSupportLang.Add(LangTh)
	allSupportLang.Add(LangZh)
}

// 当前上下文中是否为指定语言，参数应该为LangEn等的枚举类型
func IsLang(c context.Context, target string) bool {
	return compareIgnoreCase(GetLang(c), target)
}

// 判断给定src语言是否为指定target语言，target参数应该为LangEn等的枚举类型
func IsLangFor(src, target string) bool {
	return compareIgnoreCase(src, target)
}

// 当前上下文中是否为指定地区，参数应该为LocalMy等的枚举类型
func IsLocal(target string) bool {
	return compareIgnoreCase(localInner, target)
}

// 判断给定src地区是否为指定target地区，target参数应该为LocalMy等的枚举类型
func IsLocalFor(src, target string) bool {
	return compareIgnoreCase(src, target)
}

func compareIgnoreCase(s1, s2 string) bool {
	return strings.ToLower(s1) == strings.ToLower(s2)
}
