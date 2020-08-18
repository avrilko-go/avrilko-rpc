package util

import "unsafe"

// 字节切片转换为字符串
func SliceByteToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
