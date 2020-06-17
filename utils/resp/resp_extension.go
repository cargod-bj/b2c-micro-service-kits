package utils

import (
	"fmt"
	"github.com/cargod-bj/b2c-common/utils"
	commonUtils "github.com/cargod-bj/b2c-common/utils"
	"github.com/cargod-bj/b2c-ms-kits/log"
	"github.com/cargod-bj/b2c-proto-common/common"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/types/known/anypb"
	"reflect"
)

import "github.com/cargod-bj/b2c-common/resp"

// 对Response对象做的扩展方法，便于数据处理
type ResponseExtension interface {

	// 初始化当前response的code和msg为成功
	InitSuccess()

	// 解析response中的data到指定out中
	// 有可能报错，参见：resp.FAILED_DTO_DATA_NIL、resp.FAILED_DTO_DECODE
	// 返回结果(否有错误)：true 处理失败，false 处理成功
	ParseData2Dto(out interface{}) bool

	// 解析微服务入参中的dto到po中
	// 有可能报错，参见：resp.FAILED_DTO_DECODE
	// 返回结果(否有错误)：true 处理失败，false 处理成功
	ParseParams2Dto(in, out interface{}) bool

	HoldError(err error) bool

	GetFormatMsg() string

	// 验证当前response是否已经发生错误
	IsError() bool
}

func ParseData2Dto(x *common.Response, out interface{}) bool {
	data := x.Data
	if data == nil {
		x.Code = resp.FAILED_DTO_DATA_NIL
		x.Msg = resp.FAILED_DTO_DATA_NIL_MSG
		return true
	}
	if err := utils.DecodeDto(data, out); err != nil {
		x.Code = resp.FAILED_DTO_DECODE
		x.Msg = resp.FAILED_DTO_DECODE_MSG
		log.Info("解析数据错误", x, err)
		return true
	}

	return false
}

func ParseParams2Dto(x *common.Response, in, out interface{}) bool {
	data := in
	if data == nil {
		x.Code = resp.FAILED_DTO_DECODE
		x.Msg = resp.FAILED_DTO_DECODE_MSG
		return true
	}
	if err := utils.DecodeDto(data, out); err != nil {
		x.Code = resp.FAILED_DTO_DECODE
		x.Msg = resp.FAILED_DTO_DECODE_MSG
		log.Info("解析数据错误", x, err)
		return true
	}
	return false
}

func HoldError(x *common.Response, err error) bool {
	if err != nil {
		x.Code = resp.FAILED_UNKNOWN
		x.Msg = resp.FAILED_UNKNOWN_MSG
		log.Info("数据处理错误", err)
		return true
	}
	return false
}

func IsError(x *common.Response) bool {
	return x.Code != "" && x.Code != resp.SUCCESS
}

func InitSuccess(x *common.Response) {
	x.Code = resp.SUCCESS
	x.Msg = resp.SUCCESS_MSG
}

// 将DTO中的信息绑定到target上，并将其转换为any类型。
// 如果在转换过程中出现任何错误，将终止，并在response上记录错误code和msg
func Dto2Any(x *common.Response, src interface{}, target proto.Message) (result *anypb.Any, err error) {
	if err = commonUtils.DecodeDto(src, target); err != nil {
		x.Code = resp.FAILED_DTO_DECODE
		x.Msg = resp.FAILED_DTO_DECODE_MSG
		log.Info(resp.FAILED_DTO_DECODE_MSG, err)
		return
	}
	result, err = ptypes.MarshalAny(target)
	if err != nil {
		x.Code = resp.FAILED_DTO_DECODE
		x.Msg = resp.FAILED_DTO_DECODE_MSG
		log.Info(resp.FAILED_DTO_DECODE_MSG, err)
	}
	return
}

// src必须为slice类型，将src转化为target类型的slice，最终将其转换为any类型的slice。
// 如果在转换过程中出现任何错误，将终止，并在response上记录错误code和msg
func DtoList2AnyList(x *common.Response, src interface{}, target proto.Message) (resultList []*anypb.Any, err error) {
	srcType := reflect.TypeOf(src)
	if srcType.Kind() != reflect.Slice {
		x.Code = resp.FAILED_DTO_DECODE
		x.Msg = resp.FAILED_DTO_DECODE_MSG
		log.Info(resp.FAILED_DTO_DECODE_MSG, err)
		return
	}

	srcValue := reflect.ValueOf(src)
	targetType := reflect.TypeOf(target).Elem()

	srcLen := srcValue.Len()

	for i := 0; i < srcLen; i++ {
		itemValue := srcValue.Index(i)
		value := itemValue.Interface()
		if value == nil {
			break
		}
		targetInstance := reflect.New(targetType).Interface()
		if err = commonUtils.DecodeDto(value, targetInstance); err != nil {
			x.Code = resp.FAILED_DTO_DECODE
			x.Msg = resp.FAILED_DTO_DECODE_MSG
			log.Info(resp.FAILED_DTO_DECODE_MSG, err)
			return
		}
		itemAny, err := ptypes.MarshalAny(targetInstance.(proto.Message))
		if err != nil {
			x.Code = resp.FAILED_DTO_DECODE
			x.Msg = resp.FAILED_DTO_DECODE_MSG
			log.Info(resp.FAILED_DTO_DECODE_MSG, err)
		}
		resultList = append(resultList, itemAny)
	}

	return
}

// 将DTO中的信息转换到target上，再将target转为any类型，最后绑定any到data上。
// 如果你要处理列表数据，请参考使用 BindDtoList2DataList 方法。本方法只处理src为单数形式
// 如果在转换过程中出现任何错误，将终止，并在response上记录错误code和msg
func BindDto2Data(x *common.Response, src interface{}, target proto.Message) {
	if err := commonUtils.DecodeDto(src, target); err != nil {
		x.Code = resp.FAILED_DTO_DECODE
		x.Msg = resp.FAILED_DTO_DECODE_MSG
		log.Info(resp.FAILED_DTO_DECODE_MSG, err)
		return
	}
	result, err := ptypes.MarshalAny(target)
	if err != nil {
		x.Code = resp.FAILED_DTO_DECODE
		x.Msg = resp.FAILED_DTO_DECODE_MSG
		log.Info(resp.FAILED_DTO_DECODE_MSG, err)
	}
	x.Data = result
	return
}

// 此方法将slice类型的src最终转为response上的Data。
// src为数据源，一般为数据库获取的列表，必须为slice类型，否则将报错。如果src为单数，应该使用 BindDto2Data 方法。
// pageList参数将接收src转换后的[]*any，并将其赋值给response.Data。
// 所以在调用本方法前，确保pageList其他参数已经设置完成。之后的pageList更改，将不会作用到最终的response上。
// 如果在转换过程中出现任何错误，将终止，并在response上记录错误code和msg
func BindDtoList2DataList(x *common.Response, src interface{}, target proto.Message, pageList *common.PagedList) {
	srcType := reflect.TypeOf(src)
	if srcType.Kind() == reflect.Slice {
		list, err := DtoList2AnyList(x, src, target)
		if err != nil || IsError(x) {
			return
		}
		pageList.List = list
		result, err := ptypes.MarshalAny(pageList)
		if err != nil {
			x.Code = resp.FAILED_DTO_DECODE
			x.Msg = resp.FAILED_DTO_DECODE_MSG
			log.Info(resp.FAILED_DTO_DECODE_MSG, err)
			return
		}
		x.Data = result
	} else {
		x.Code = resp.FAILED_DTO_DECODE
		x.Msg = resp.FAILED_DTO_DECODE_MSG
		log.Info(resp.FAILED_DTO_DECODE_MSG)
	}
	return
}

// 此方法将Message类型的slice最终转为response上的Data。
// list为数据源，必须为Message类型的slice。
// pageList参数将接收src转换后的[]*any，并将其赋值给response.Data。
// 所以在调用本方法前，确保pageList其他参数已经设置完成。之后的pageList更改，将不会作用到最终的response上。
// 如果在转换过程中出现任何错误，将终止，并在response上记录错误code和msg
func BindMessageList(x *common.Response, list interface{}, pageList *common.PagedList) {
	defer func() {
		if err := recover(); err != nil {
			log.Info(resp.FAILED_DTO_DECODE_MSG, err)
		}
	}()
	srcType := reflect.TypeOf(list)
	if srcType.Kind() == reflect.Slice {
		listValue := reflect.ValueOf(list)
		if listValue.Kind() == reflect.Ptr {
			listValue = listValue.Elem()
		}
		var anyList []*anypb.Any
		for i := 0; i < listValue.Len(); i++ {
			var item interface{}
			itemValue := listValue.Index(i)
			if itemValue.Kind() == reflect.Ptr {
				item = itemValue.Interface()
			} else {
				item = itemValue.Interface()
			}
			itemAny, err := ptypes.MarshalAny(item.(proto.Message))
			if err != nil {
				x.Code = resp.FAILED_DTO_DECODE
				x.Msg = resp.FAILED_DTO_DECODE_MSG
				log.Info(resp.FAILED_DTO_DECODE_MSG, err)
				return
			}
			anyList = append(anyList, itemAny)
		}
		pageList.List = anyList
		result, err := ptypes.MarshalAny(pageList)
		if err != nil {
			x.Code = resp.FAILED_DTO_DECODE
			x.Msg = resp.FAILED_DTO_DECODE_MSG
			log.Info(resp.FAILED_DTO_DECODE_MSG, err)
			return
		}
		x.Data = result
	} else {
		x.Code = resp.FAILED_DTO_DECODE
		x.Msg = resp.FAILED_DTO_DECODE_MSG
		log.Info(resp.FAILED_DTO_DECODE_MSG)
	}
	return
}

func GetFormatMsg(x *common.Response) string {
	return fmt.Sprintf("%s(%s)", x.Msg, x.Code)
}
