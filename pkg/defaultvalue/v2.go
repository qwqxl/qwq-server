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
//	// 处理指针：如果是nil则初始化，然后递归处理指向的值
//	if v.Kind() == reflect.Ptr {
//		if v.IsNil() {
//			v.Set(reflect.New(v.Type().Elem()))
//		}
//		return setDefaults(v.Elem(), checkTag)
//	}
//
//	// 只处理结构体和切片
//	switch v.Kind() {
//	case reflect.Struct:
//		// 处理结构体字段
//		for i := 0; i < v.NumField(); i++ {
//			field := v.Field(i)
//			fieldType := v.Type().Field(i)
//
//			// 跳过不可导出字段
//			if !field.CanSet() {
//				continue
//			}
//
//			tagValue := fieldType.Tag.Get(DefaultValTagName)
//
//			// 递归处理字段
//			if err := setDefaults(field, true); err != nil {
//				return err
//			}
//
//			// 处理切片字段
//			if field.Kind() == reflect.Slice {
//				// 有标签且切片为空
//				if checkTag && tagValue != "" && (field.IsNil() || field.Len() == 0) {
//					length, err := strconv.Atoi(tagValue)
//					if err != nil {
//						return fmt.Errorf("field %s: slice length must be integer: %w", fieldType.Name, err)
//					}
//					if length < 0 {
//						return fmt.Errorf("field %s: slice length must be non-negative", fieldType.Name)
//					}
//
//					// 创建新切片
//					sliceType := reflect.SliceOf(field.Type().Elem())
//					newSlice := reflect.MakeSlice(sliceType, length, length)
//					field.Set(newSlice)
//				}
//
//				// 递归处理切片元素
//				for j := 0; j < field.Len(); j++ {
//					if err := setDefaults(field.Index(j), true); err != nil {
//						return fmt.Errorf("field %s: element %d: %w", fieldType.Name, j, err)
//					}
//				}
//				continue
//			}
//
//			// 处理基本类型标签
//			if checkTag && tagValue != "" && isZero(field) {
//				if err := setFieldValue(field, tagValue); err != nil {
//					return fmt.Errorf("field %s: %w", fieldType.Name, err)
//				}
//			}
//		}
//
//	case reflect.Slice:
//		// 递归处理切片元素
//		for i := 0; i < v.Len(); i++ {
//			if err := setDefaults(v.Index(i), true); err != nil {
//				return err
//			}
//		}
//	}
//
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
//	case reflect.Ptr:
//		return v.IsNil() || isZero(v.Elem())
//	case reflect.Struct:
//		for i := 0; i < v.NumField(); i++ {
//			if !isZero(v.Field(i)) {
//				return false
//			}
//		}
//		return true
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
//	case reflect.Slice, reflect.Map, reflect.Interface, reflect.Chan, reflect.Func:
//		return v.IsNil() || v.Len() == 0
//	case reflect.Array:
//		for i := 0; i < v.Len(); i++ {
//			if !isZero(v.Index(i)) {
//				return false
//			}
//		}
//		return true
//	default:
//		return false
//	}
//}
