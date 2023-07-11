package convert

import (
	"fmt"
	"unsafe"
)

var convertList = []string{"B", "KB", "MB", "GB", "TB", "PB"}

const base = float64(1024)

func Bytes2Human(n int) string {
	index := 0
	floatn := float64(n)
	for floatn >= base {
		floatn /= base
		index++
	}
	return fmt.Sprintf("%.2f%s", floatn, convertList[index])
}

func Str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
