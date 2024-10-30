package handler

import (
	"encoding/json" // JSON シリアライズのために追加
	"log"           // エラーログ出力のために追加
	"net/http"

	"github.com/TechBowl-japan/go-stations/model"
)

// HealthzHandler はヘルスチェックエンドポイントのハンドラーです。
type HealthzHandler struct{}

// NewHealthzHandler は HealthzHandler を返します。
func NewHealthzHandler() *HealthzHandler {
	return &HealthzHandler{}
}

// ServeHTTP は http.Handler インターフェースを実装します。
func (h *HealthzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// レスポンスヘッダーに Content-Type を設定
	w.Header().Set("Content-Type", "application/json")

	// ステータスコードを 200 OK に設定
	w.WriteHeader(http.StatusOK)

	// HealthzResponse 構造体のインスタンスを作成し、Message に "OK" を設定
	response := model.HealthzResponse{
		Message: "OK",
	}

	// JSON エンコードしてレスポンスとして送信
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		// エンコードに失敗した場合はエラーログを出力
		log.Println("JSON Encode Error:", err)
		return
	}
}
