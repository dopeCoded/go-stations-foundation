//station2
package middleware

import (
    "context"
    "log"
    "net/http"

    "github.com/mileusna/useragent"
)

// contextKey は Context のキーとして使用するカスタム型です。
type contextKey string

const OSContextKey contextKey = "os"

// OSExtractor は、User-AgentからOS名を抽出してContextに格納するミドルウェアです。
func OSExtractor(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // User-Agent ヘッダーを取得
        uaString := r.UserAgent()

        // User-Agent を解析
        ua := useragent.Parse(uaString)

        // OS 名を取得
        osName := ua.OS
        if osName == "" {
            osName = "Unknown"
        }

        // デバッグ用ログ出力
        log.Printf("Middleware: Parsed OS as '%s' from User-Agent '%s'", osName, uaString)

        // Context に OS 名を格納
        ctx := context.WithValue(r.Context(), OSContextKey, osName)

        // 新しい Context を持つ Request を作成
        r = r.WithContext(ctx)

        // 次のハンドラーを呼び出す
        next.ServeHTTP(w, r)
    })
}

//station2 end
