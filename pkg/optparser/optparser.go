package optparser

import "strings"

func ParseOpts(raw string, ignoreEmpty bool) ([]string, map[string]string) {
	var positional []string
	var keyValue map[string]string

	for _, str := range strings.Split(raw, ",") {
		if str == "" && ignoreEmpty {
			continue
		}
		if idx := strings.Index(str, "="); idx != -1 {
			if keyValue == nil {
				keyValue = make(map[string]string)
			}
			keyValue[str[:idx]] = str[idx+1:]
		} else {
			positional = append(positional, str)
		}
	}

	return positional, keyValue
}
