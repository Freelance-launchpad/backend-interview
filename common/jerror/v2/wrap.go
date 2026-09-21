package jerror

import (
	"fmt"
	"runtime"
	"strings"
)

func Wrap(errPtr *error, params ...string) {
	if errPtr != nil && *errPtr != nil {
		pc := make([]uintptr, 10)
		runtime.Callers(2, pc)
		f := runtime.FuncForPC(pc[0])

		parts := strings.Split(f.Name(), "/")
		funcName := strings.Join(parts[len(parts)-1:], "")

		var s string
		if len(params) > 0 {
			s = " " + strings.Join(params, " ")
		}

		*errPtr = fmt.Errorf("%s has failed%s: %w", funcName, s, *errPtr)
	}
}
