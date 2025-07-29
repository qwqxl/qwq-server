package defaultvalue

//
//import (
//	"fmt"
//	"reflect"
//	"strconv"
//	"time"
//)
//
//const (
//	// DefaultValTagName 标签名
//	DefaultValTagName = "qwq-default"
//)
//
//// SetDefaults 设置结构体的默认值（接受结构体指针）
//func SetDefaults(ptr interface{}) error {
//	v := reflect.ValueOf(ptr)
//	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
//		return fmt.Errorf("input must be a pointer to struct")
//	}
//	return setDefaults(v.Elem(), true)
//}
//
//// 递归设置默认值的核心函数
//func setDefaults(v reflect.Value, checkTag bool) error {
//	//t := v.Type()
//
//	// 如果是nil指针，则创建新实例
//	if v.Kind() == reflect.Ptr && v.IsNil() && v.Type().Elem().Kind() == reflect.Struct {
//		v.Set(reflect.New(v.Type().Elem()))
//	}
//
//	// 获取实际值（如果是指针则解引用）
//	actualValue := v
//	if v.Kind() == reflect.Ptr {
//		actualValue = v.Elem()
//	}
//
//	// 只处理结构体
//	if actualValue.Kind() != reflect.Struct {
//		return nil
//	}
//
//	for i := 0; i < actualValue.NumField(); i++ {
//		field := actualValue.Field(i)
//		fieldType := actualValue.Type().Field(i)
//
//		// 跳过不可导出字段
//		if !field.CanSet() {
//			continue
//		}
//
//		tagValue := fieldType.Tag.Get(DefaultValTagName)
//		fieldKind := field.Kind()
//
//		// 处理嵌套结构体
//		if fieldKind == reflect.Struct ||
//			(fieldKind == reflect.Ptr && field.Type().Elem().Kind() == reflect.Struct) {
//
//			// 确保指针已初始化
//			if fieldKind == reflect.Ptr && field.IsNil() {
//				field.Set(reflect.New(field.Type().Elem()))
//			}
//
//			// 递归处理嵌套结构体
//			target := field
//			if fieldKind == reflect.Ptr {
//				target = field.Elem()
//			}
//
//			if err := setDefaults(target, true); err != nil {
//				return err
//			}
//			continue
//		}
//
//		// 处理带标签的基本类型
//		if tagValue != "" && isZero(field) {
//			if err := setFieldValue(field, tagValue); err != nil {
//				return fmt.Errorf("field %s: %w", fieldType.Name, err)
//			}
//		}
//	}
//	return nil
//}
//
//// 设置字段值
//func setFieldValue(field reflect.Value, value string) error {
//	// 处理time.Duration类型（包括指针和非指针）
//	if field.Type() == reflect.TypeOf(time.Duration(0)) {
//		dur, err := time.ParseDuration(value)
//		if err != nil {
//			return fmt.Errorf("invalid duration: %w", err)
//		}
//		field.Set(reflect.ValueOf(dur))
//		return nil
//	}
//
//	if field.Type() == reflect.TypeOf((*time.Duration)(nil)).Elem() {
//		dur, err := time.ParseDuration(value)
//		if err != nil {
//			return fmt.Errorf("invalid duration: %w", err)
//		}
//		if field.IsNil() {
//			field.Set(reflect.ValueOf(&dur))
//		} else {
//			field.Elem().Set(reflect.ValueOf(dur))
//		}
//		return nil
//	}
//
//	// 处理其他类型
//	switch field.Kind() {
//	case reflect.String:
//		field.SetString(value)
//
//	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
//		val, err := strconv.ParseInt(value, 10, 64)
//		if err != nil {
//			return err
//		}
//		field.SetInt(val)
//
//	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
//		val, err := strconv.ParseUint(value, 10, 64)
//		if err != nil {
//			return err
//		}
//		field.SetUint(val)
//
//	case reflect.Float32, reflect.Float64:
//		val, err := strconv.ParseFloat(value, 64)
//		if err != nil {
//			return err
//		}
//		field.SetFloat(val)
//
//	case reflect.Bool:
//		val, err := strconv.ParseBool(value)
//		if err != nil {
//			return err
//		}
//		field.SetBool(val)
//
//	case reflect.Ptr: // 处理基本类型的指针
//		if field.IsNil() {
//			field.Set(reflect.New(field.Type().Elem()))
//		}
//		return setFieldValue(field.Elem(), value)
//
//	default:
//		return fmt.Errorf("unsupported type: %s", field.Type())
//	}
//	return nil
//}
//
//// 检查是否为零值
//func isZero(v reflect.Value) bool {
//	switch v.Kind() {
//	case reflect.String:
//		return v.Len() == 0
//	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
//		return v.Int() == 0
//	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
//		return v.Uint() == 0
//	case reflect.Float32, reflect.Float64:
//		return v.Float() == 0
//	case reflect.Bool:
//		return !v.Bool()
//	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
//		return v.IsNil()
//	}
//	return false
//}
