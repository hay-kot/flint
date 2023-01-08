package flint

import (
	"errors"
	"fmt"
)

type FlintErrors map[string][]error

type FileError struct {
	Path string
	Err  error
}

func (fe FileError) Error() string {
	return fmt.Sprintf("%s: %s", fe.Path, fe.Err)
}

func IsFileError(err error) bool {
	var e FileError
	return errors.As(err, &e)
}

func (fe FlintErrors) Error() string {
	return "flint errors"
}
