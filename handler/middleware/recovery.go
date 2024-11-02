//station1
package middleware

import (
    "log"
    "net/http"
)

// Recovery はパニックをキャッチしてアプリケーションのクラッシュを防ぐミドルウェアです。
func Recovery(h http.Handler) http.Handler {
    fn := func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                // エラーログを出力
                log.Printf("panic recovered: %v", err)
                // HTTP 500 Internal Server Error を返す
                http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
            }
        }()

        // 次のハンドラーを呼び出す
        h.ServeHTTP(w, r)
    }
    return http.HandlerFunc(fn)
}
//station1 end
