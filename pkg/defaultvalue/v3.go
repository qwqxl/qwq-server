package defaultvalue

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	// DefaultValTagName 标签名
	DefaultValTagName = "qwq-default"
)

// SetDefaults 设置结构体的默认值（接受结构体指针）
func SetDefaults(ptr interface{}) error {
	v := reflect.ValueOf(ptr)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("input must be a pointer to struct")
	}
	return setDefaults(v.Elem(), true)
}

// 递归设置默认值的核心函数
func setDefaults(v reflect.Value, checkTag bool) error {
	// 处理指针：如果是nil则初始化，然后递归处理指向的值
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		return setDefaults(v.Elem(), checkTag)
	}

	// 只处理结构体和切片
	switch v.Kind() {
	case reflect.Struct:
		// 处理结构体字段
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			fieldType := v.Type().Field(i)

			// 跳过不可导出字段
			if !field.CanSet() {
				continue
			}

			tagValue := fieldType.Tag.Get(DefaultValTagName)

			// 处理切片字段
			if field.Kind() == reflect.Slice {
				// 有标签且切片为空
				if checkTag && tagValue != "" && (field.IsNil() || field.Len() == 0) {
					if err := initSliceFromTag(field, tagValue); err != nil {
						return fmt.Errorf("field %s: %w", fieldType.Name, err)
					}
				}
			}

			// 递归处理字段
			if err := setDefaults(field, true); err != nil {
				return err
			}

			// 处理基本类型标签
			if checkTag && tagValue != "" && isZero(field) {
				if err := setFieldValue(field, tagValue); err != nil {
					return fmt.Errorf("field %s: %w", fieldType.Name, err)
				}
			}
		}

	case reflect.Slice:
		// 递归处理切片元素
		for i := 0; i < v.Len(); i++ {
			if err := setDefaults(v.Index(i), true); err != nil {
				return err
			}
		}
	}

	return nil
}

// 从标签初始化切片
func initSliceFromTag(slice reflect.Value, tagValue string) error {
	//elemType := slice.Type().Elem()

	// 尝试解析为JSON数组
	if strings.HasPrefix(tagValue, "[") && strings.HasSuffix(tagValue, "]") {
		return initSliceFromJSON(slice, tagValue)
	}

	// 尝试解析为逗号分隔的值
	if strings.Contains(tagValue, ",") {
		return initSliceFromCSV(slice, tagValue)
	}

	// 否则尝试解析为单个值
	return initSliceFromSingleValue(slice, tagValue)
}

// 从JSON数组初始化切片
func initSliceFromJSON(slice reflect.Value, jsonStr string) error {
	//elemType := slice.Type().Elem()

	// 创建目标切片的指针
	slicePtr := reflect.New(slice.Type())
	if err := json.Unmarshal([]byte(jsonStr), slicePtr.Interface()); err != nil {
		return fmt.Errorf("invalid JSON array: %w", err)
	}

	slice.Set(slicePtr.Elem())
	return nil
}

// 从逗号分隔的值初始化切片
func initSliceFromCSV(slice reflect.Value, csv string) error {
	//elemType := slice.Type().Elem()
	values := strings.Split(csv, ",")

	// 创建新切片
	newSlice := reflect.MakeSlice(slice.Type(), len(values), len(values))

	for i, valStr := range values {
		valStr = strings.TrimSpace(valStr)
		if err := setFieldValue(newSlice.Index(i), valStr); err != nil {
			return fmt.Errorf("element %d: %w", i, err)
		}
	}

	slice.Set(newSlice)
	return nil
}

// 从单个值初始化切片（创建长度为1的切片）
func initSliceFromSingleValue(slice reflect.Value, value string) error {
	// 尝试解析为整数（长度）
	if length, err := strconv.Atoi(value); err == nil && length >= 0 {
		newSlice := reflect.MakeSlice(slice.Type(), length, length)
		slice.Set(newSlice)
		return nil
	}

	// 否则创建长度为1的切片并设置值
	newSlice := reflect.MakeSlice(slice.Type(), 1, 1)
	if err := setFieldValue(newSlice.Index(0), value); err != nil {
		return err
	}

	slice.Set(newSlice)
	return nil
}

// 设置字段值
func setFieldValue(field reflect.Value, value string) error {
	// 处理time.Duration类型
	if field.Type() == reflect.TypeOf(time.Duration(0)) {
		dur, err := time.ParseDuration(value)
		if err != nil {
			return fmt.Errorf("invalid duration: %w", err)
		}
		field.Set(reflect.ValueOf(dur))
		return nil
	}

	if field.Type() == reflect.TypeOf((*time.Duration)(nil)).Elem() {
		dur, err := time.ParseDuration(value)
		if err != nil {
			return fmt.Errorf("invalid duration: %w", err)
		}
		if field.IsNil() {
			field.Set(reflect.ValueOf(&dur))
		} else {
			field.Elem().Set(reflect.ValueOf(dur))
		}
		return nil
	}

	// 处理time.Time类型
	if field.Type() == reflect.TypeOf(time.Time{}) {
		// 支持常见时间格式
		formats := []string{
			time.RFC3339,
			"2006-01-02",
			"2006-01-02 15:04:05",
		}

		var t time.Time
		var err error
		for _, format := range formats {
			t, err = time.Parse(format, value)
			if err == nil {
				break
			}
		}

		if err != nil {
			return fmt.Errorf("invalid time format: %w", err)
		}

		field.Set(reflect.ValueOf(t))
		return nil
	}

	// 处理其他类型
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(val)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(val)

	case reflect.Float32, reflect.Float64:
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(val)

	case reflect.Bool:
		val, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(val)

	case reflect.Ptr: // 处理基本类型的指针
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		return setFieldValue(field.Elem(), value)

	default:
		return fmt.Errorf("unsupported type: %s", field.Type())
	}
	return nil
}

// 检查是否为零值
func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Ptr:
		return v.IsNil() || isZero(v.Elem())
	case reflect.Struct:
		// 特殊处理time.Time
		if t, ok := v.Interface().(time.Time); ok {
			return t.IsZero()
		}
		for i := 0; i < v.NumField(); i++ {
			if !isZero(v.Field(i)) {
				return false
			}
		}
		return true
	case reflect.String:
		return v.Len() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Slice, reflect.Map, reflect.Interface, reflect.Chan, reflect.Func:
		return v.IsNil() || v.Len() == 0
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if !isZero(v.Index(i)) {
				return false
			}
		}
		return true
	default:
		return false
	}
}
