package strings

import (
	"strings"
)

func SplitString(str string, sep string) []string {
	if string(str) == "" {
		return nil
	}
	return strings.Split(str, sep)
}

func SplitStringToMapBool(val, sep string) map[string]bool {
	if val == "" {
		return nil
	}
	ret := make(map[string]bool)
	for _, v := range strings.Split(val, sep) {
		ret[v] = true
	}
	return ret
}

func InSeparatedString(val, check, sep string) bool {
	return SplitStringToMapBool(val, sep)[check]
}

func AppendSeparatedString(val, sep string, values ...string) string {
	nv := strings.Join(values, sep)
	if val == "" {
		return nv
	}
	return val + sep + nv
}
