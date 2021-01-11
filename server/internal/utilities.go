package internal

import (
	"fmt"
)

// FailOnError reduces boilerplate
func FailOnError(err error, msg string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
