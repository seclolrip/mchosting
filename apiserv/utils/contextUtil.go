package utils

import (
	"runtime"
	"strings"
)

func GetCurrentFuncName() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return "unknown"
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}

	fullName := fn.Name()
	parts := strings.Split(fullName, "/")
	shortName := parts[len(parts)-1]

	return shortName
}
