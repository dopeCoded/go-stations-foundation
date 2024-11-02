//station4
package middleware

import (
	"net/http"
)

type BasicAuthMiddleware struct {
	UserID   string
	Password string
}

// NewBasicAuthMiddleware は Basic 認証ミドルウェアを作成します。
func NewBasicAuthMiddleware(userID, password string) *BasicAuthMiddleware {
	return &BasicAuthMiddleware{
		UserID:   userID,
		Password: password,
	}
}

// Handler は Basic 認証を行うハンドラーを返します。
func (bam *BasicAuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストからユーザー名とパスワードを取得
		userID, password, ok := r.BasicAuth()
		if !ok {
			// 認証情報がない場合
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// 環境変数から取得したユーザー名とパスワードと比較
		if userID != bam.UserID || password != bam.Password {
			// 認証失敗
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// 認証成功、次のハンドラーを呼び出す
		next.ServeHTTP(w, r)
	})
}
