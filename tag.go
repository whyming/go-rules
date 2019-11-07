package gorules

import (
	"reflect"
	"unsafe"
)

const (
	sStart = iota
	sFound
	sNotFound
)

func getTagName(t reflect.StructTag) string {
	s := *(*string)(unsafe.Pointer(&t))

	var start, l, status int
Tag:
	for i := 0; i < len(s); i++ {
		switch status {
		case sStart:
			if s[i] != ' ' {
				var end int
				for i+end < len(s) && s[i+end] != ':' {
					end++
				}
				if s[i:i+end] == "json" {
					status = sFound
				} else {
					status = sNotFound
				}
				i += end
			}

		case sFound:
			if s[i] != '"' {
				break Tag
			}
			start = i + 1
			for start < len(s) && s[i] == ' ' {
				start++
			}

			for start+l < len(s) && s[start+l] != ',' && s[start+l] != '"' {
				l++
			}
			break Tag
		case sNotFound:
			i++
			for i < len(s) && s[i] != '"' {
				i++
			}
			status = sStart
		}
	}
	return s[start : start+l]
}
