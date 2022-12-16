package utilz

import "errors"

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

func IsInArray[T BasicType](resultStatus T, array []T) bool {
	for _, status := range array {
		if status == resultStatus {
			return true
		}
	}
	return false
}
func IsInArray2[T int | int64 | int32 | int8 | int16 | float32 | float64 | bool | string](resultStatus T, array []T) bool {
	for _, status := range array {
		if status == resultStatus {
			return true
		}
	}
	return false
}
func MapMerge(mp ...map[string]string) map[string]string {
	m2 := make(map[string]string)
	for _, m1 := range mp {
		for k, v := range m1 {
			m2[k] = v
		}
	}
	return m2
}
func GetMapValue(v map[string]any, path ...string) (any, error) {
	if len(path) == 0 {
		return v, nil
	}
	ok := false
	for i := 0; i < len(path)-1; i++ {
		v, ok = v[path[i]].(map[string]any)
		if !ok {
			return nil, errors.New(path[i] + " fail")
		}
	}
	return v[path[len(path)-1]], nil
}
