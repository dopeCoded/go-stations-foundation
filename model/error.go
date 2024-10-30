package model

// ErrNotFound は、対象の TODO が存在しない場合のエラーを表します。
type ErrNotFound struct{}

// Error メソッドを実装し、error インターフェースを満たします。
func (e ErrNotFound) Error() string {
	return "the requested resource was not found"
}