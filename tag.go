package gorules

import (
	"reflect"
)

const (
	sStart = iota
	sFound
	sNotFound
)

// getTagName if has rule use rule tag ,then json
func getTagName(t reflect.StructTag) string {
	name := t.Get("rule")
	if name == "" {
		name = t.Get("json")
	}

	for i, s := range name {
		if s == ',' {
			return name[0:i]
		}
	}
	return name
}
