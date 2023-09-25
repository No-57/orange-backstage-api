package convert

import "unsafe"

func BytesToStr(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func StrToBytes(s string) (b []byte) {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
