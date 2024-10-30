package model

import (
	"errors"
)

// ErrNotFound は、対象の TODO が存在しない場合のエラーを表します。
type ErrNotFound struct{}

// Error メソッドを実装し、error インターフェースを満たします。
func (e ErrNotFound) Error() string {
	return "the requested resource was not found"
}

func IsErrNotFound(err error) bool {
	var notFoundErr ErrNotFound
	return errors.As(err, &notFoundErr)
}