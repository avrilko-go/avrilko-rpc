package util

import "unsafe"

// 字节切片转换为字符串
func SliceByteToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// 字符串转字节切片
func StringToByteSlice(s string) []byte {
	sPtr := (*[2]uintptr)(unsafe.Pointer(&s))     // 字符串的指针底层结构（ []uintptr{ptr,len} ）
	bPtr := [3]uintptr{sPtr[0], sPtr[1], sPtr[1]} // 字符切片的指针底层结构（ []uintptr{ptr,len,cap} ）

	return *(*[]byte)(unsafe.Pointer(&bPtr))
}
