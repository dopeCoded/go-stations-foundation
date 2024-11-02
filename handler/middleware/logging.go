// station3
package middleware

import (
	"fmt"
	"net/http"
	"time"
)

// LogEntry represents the structure of the log to be output.
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Latency   int64     `json:"latency"` // in milliseconds
	Path      string    `json:"path"`
	OS        string    `json:"os"`
}

// LoggingMiddleware logs the request information before and after handler execution.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. リクエストの前処理
		startTime := time.Now()

		// 次のハンドラーを呼び出す
		next.ServeHTTP(w, r)

		// 2. リクエストの後処理
		latency := time.Since(startTime).Milliseconds()

		// 3. Context から OS 情報を取得
		osName, ok := r.Context().Value(OSContextKey).(string)
		if !ok || osName == "" {
			osName = "Unknown"
		}

		// 4. URL パスを取得
		path := r.URL.Path

		// 5. LogEntry を作成
		logEntry := LogEntry{
			Timestamp: startTime,
			Latency:   latency,
			Path:      path,
			OS:        osName,
		}

		// 6. JSON にシリアライズして標準出力に書き出し
		// logJSON, err := json.Marshal(logEntry)
		// if err != nil {
		//     log.Printf("Error marshaling log entry: %v", err)
		//     return
		// }
		fmt.Printf("%+v", logEntry)
	})
}

//station3 end
