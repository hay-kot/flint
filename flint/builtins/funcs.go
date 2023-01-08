package builtins

import "fmt"

// mapAnyToStr maps any slice to string slice and calls f for each string.
// values that don't pass a type assertion to string are ignored.
func mapAnyToStr(vv []any, f func(s string, i int)) {
	for i, v := range vv {
		s, ok := v.(string)
		if !ok {
			continue
		}

		f(s, i)
	}
}

func fmtKeyCords(x, y int) string {
	if x == -1 {
		return "0:0"
	}

	return fmt.Sprintf("%d:%d", x, y)
}
