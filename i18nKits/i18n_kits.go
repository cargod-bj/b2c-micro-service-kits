package i18nKits

import (
	"context"
	"fmt"
	"github.com/cargod-bj/b2c-common/commonUtils/goroutineKits/gBox"
	"github.com/cargod-bj/b2c-proto-common/common"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/server"
	"strings"
)

// 通过context传向微服务的i18n内容的key值
const I18nContextKey = "language$Value$Key"

type LogFun = func(...interface{})

var (
	logErr  LogFun
	logInfo LogFun
)

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

// 添加i18n的handler
// 使用方式如下：
// 	var service micro.Service
//	初始化service
//	service.Init(i18nKits.I18nMicroHandler)
// 本方法只添加了i18n信息
func I18nMicroHandler(opt *micro.Options) {
	_ = opt.Server.Init(func(o *server.Options) {
		o.HdlrWrappers = append(o.HdlrWrappers, func(originalHandler server.HandlerFunc) server.HandlerFunc {
			return func(ctx context.Context, req server.Request, rsp interface{}) error {
				defer func() {
					if r := recover(); r != nil {
						//if logErr != nil {
						//	logErr("error occur:", r)
						//}
						fmt.Printf("error occur: %v", r)
						if t, ok := rsp.(*common.Response); ok {
							t.Code = defaultUnknown
						}
					}
				}()
				err := originalHandler(ctx, req, rsp)
				if t, ok := rsp.(*common.Response); ok {
					// 如果没有手动设置message，则使用code查找对应message
					if t.Msg == "" {
						t.Msg = GetRM(ctx, t.Code)
					}
				}
				if err != nil {
					fmt.Printf("error: %v", err)
				}
				//if err != nil && logErr != nil {
				//	logErr("error:", err)
				//}
				return err
			}
		})
	})
}

// 初始化log信息
// el error日志级别
// il info日志级别
func InitLog(el, il LogFun) {
	logErr = el
	logInfo = il
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

// 获取当前的语言
func GetLangByHttpHeader(acceptLanguageHeader string) string {
	if acceptLanguageHeader == "" {
		return getDefaultLangByLocal()
	}
	switch {
	case strings.Contains(acceptLanguageHeader, LangEn):
		return LangEn
	case strings.Contains(acceptLanguageHeader, LangId):
		return LangId
	case strings.Contains(acceptLanguageHeader, LangTh):
		return LangTh
	case strings.Contains(acceptLanguageHeader, LangZh):
		return LangZh
	default:
		return getDefaultLangByLocal()
	}
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
