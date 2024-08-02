package slicez

import (
	"fmt"
	"strings"
)

func JoinSlice[T any](s []T, sep any) string {
	var result string
	for _, v := range s {
		if i, ok := sep.(int32); ok {
			result += fmt.Sprintf("%v%c", v, i)
		} else {
			result += fmt.Sprintf("%v%v", v, sep)
		}
	}
	return strings.TrimSpace(result)
}

// %s 空格连接
// %s\n 换行连接
// \t%s\n 开头tab后换行连接
func JoinSliceFormat(flags []any, format string) string {
	var builder strings.Builder
	for _, flag := range flags {
		builder.WriteString(fmt.Sprintf(format, flag))
	}
	return builder.String()
}
func AppendDistinctAll(flags []string, flagsToAdd []string) []string {
	for _, s := range flagsToAdd {
		flags = AppendDistinct(flags, s)
	}
	return flags
}
func AppendDistinct(flags []string, flag string) []string {
	for _, s := range flags {
		if s == flag {
			return flags
		}
	}
	flags = append(flags, flag)
	return flags
}
