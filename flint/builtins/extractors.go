package builtins

import (
	"fmt"
)

func fmtKeyCords(x, y int) string {
	if x == -1 {
		return "0:0"
	}

	return fmt.Sprintf("%d:%d", x, y)
}
