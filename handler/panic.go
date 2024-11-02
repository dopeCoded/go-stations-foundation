//station1
package handler

import (
    "net/http"
)

// PanicHandler は常にパニックを発生させるハンドラーです。
type PanicHandler struct{}

// ServeHTTP implements http.Handler interface.
func (h *PanicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    panic("intentional panic for testing")
}
//station1 end