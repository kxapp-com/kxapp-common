package utilz

import (
	"fmt"
	"reflect"
	"strconv"
)

/*
*
系统基本的数字类型,做泛型编程，经常会要用到基本类型做类型约束，所以定义了这个基本类型
*/
type NumberType interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 |
		~uint16 | ~uint32 | ~uint64 | ~uintptr | ~float32 | ~float64
}

/*
*
系统基本类型，包括数值类型和string+bool
*/
type BasicType interface {
	NumberType | ~string | ~bool
}

func InSlice(slice any, item any) bool {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return false
		//panic("InSlice called with non-slice value")
	}
	for i := 0; i < s.Len(); i++ {
		if reflect.DeepEqual(item, s.Index(i).Interface()) {
			return true
		}
	}
	return false
}

//func IsInArray[T BasicType](resultStatus T, array []T) bool {
//	for _, status := range array {
//		if status == resultStatus {
//			return true
//		}
//	}
//	return false
//}
//func IsInArray2[T int | int64 | int32 | int8 | int16 | float32 | float64 | bool | string](resultStatus T, array []T) bool {
//	for _, status := range array {
//		if status == resultStatus {
//			return true
//		}
//	}
//	return false
//}

/*
*
从一个数组里面删除指定索引的元素
*/
func RemoveAtIndex(slice any, index int) any {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("RemoveIndex called with non-slice value")
	}
	return reflect.AppendSlice(s.Slice(0, index), s.Slice(index+1, s.Len())).Interface()
}

func MapMerge(maps ...map[any]any) map[any]any {
	result := make(map[any]any)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

/*
*
获取map或者数组内的数据，path是路径，如果是 取 array的，则path是array的index的字符串形式

	m := map[string]any{
			"foo": []any{
				map[string]any{"bar": "baz"},
				map[string]any{"bar": "qux"},
			},
		}
		val, err := GetMapArrayValue(m, "foo", "1", "bar")
*/
func GetMapValue(v any, path ...string) (any, error) {
	if len(path) == 0 {
		return v, nil
	}
	for _, p := range path {
		switch t := v.(type) {
		case map[string]any:
			var ok bool
			//v, ok = t[p.(string)]
			v, ok = t[p]
			if !ok {
				return nil, fmt.Errorf("key not found: %s", p)
			}
		case []any:
			//i, err := strconv.Atoi(p.(string))
			i, err := strconv.Atoi(p)
			if err != nil || i < 0 || i >= len(t) {
				return nil, fmt.Errorf("index out of range: %s", p)
			}
			v = t[i]
		default:
			return nil, fmt.Errorf("invalid type: %T", v)
		}
	}
	return v, nil
}
