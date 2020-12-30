package i18nKits

import (
	"context"
	"fmt"
	"github.com/cargod-bj/b2c-common/commonUtils/goroutineKits/gBox"
	"github.com/micro/go-micro/v2/metadata"
	"strings"
)

// 通过context传向微服务的i18n内容的key值
const I18nContextKey = "language$Value$Key"

// 初始化language，如果autoRecycle=true，则会自动回收，否则需要手动调用Recycle进行回收
func initI18n(acceptLanguageHeader string, autoRecycle bool) {
	switch {
	case strings.Contains(acceptLanguageHeader, LangEn):
		initInner(LangEn, autoRecycle)
	case strings.Contains(acceptLanguageHeader, LangId):
		initInner(LangId, autoRecycle)
	case strings.Contains(acceptLanguageHeader, LangTh):
		initInner(LangTh, autoRecycle)
	case strings.Contains(acceptLanguageHeader, LangZh):
		initInner(LangZh, autoRecycle)
	default:
		initInner("", autoRecycle)
	}
}

// 初始化language，如果autoRecycle=true，则会自动回收，否则需要手动调用Recycle进行回收
func initI18nByLang(language string, autoRecycle bool) {
	initInner(language, autoRecycle)
}

// 获取当前的语言
func GetLang(c context.Context) string {
	if c == nil {
		return getDefaultLangByLocal()
	}
	lang, _ := metadata.Get(c, I18nContextKey)
	if lang == "" {
		lang = getDefaultLangByLocal()
	}
	return lang
}

func initInner(language string, autoRecycle bool) {
	if language == "" {
		language = getDefaultLangByLocal()
	} else if !allSupportLang.Contain(language) {
		panic(fmt.Sprintf("暂时不支持的语言类型:%s", language))
	}
	var ar uint32
	if !autoRecycle {
		ar = 1
	}
	gBox.Put(I18nContextKey, language, ar)
}
