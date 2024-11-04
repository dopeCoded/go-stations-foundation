package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/handler/router"
)

func main() {
	err := realMain()
	if err != nil {
		log.Fatalln("main: failed to exit successfully, err =", err)
	}
}

func realMain() error {
	// config values
	const (
		defaultPort   = ":8080"
		defaultDBPath = ".sqlite3/todo.db"
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = defaultDBPath
	}

	// station4
	// 環境変数からユーザーIDとパスワードを取得
	userID := os.Getenv("BASIC_AUTH_USER_ID")
	password := os.Getenv("BASIC_AUTH_PASSWORD")

	// set time zone
	var err error
	time.Local, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return err
	}

	// set up sqlite3
	todoDB, err := db.NewDB(dbPath)
	if err != nil {
		return err
	}
	defer todoDB.Close()

	// WaitGroupを作成
	var wg sync.WaitGroup

	// station4
	// NOTE: 新しいエンドポイントの登録はrouter.NewRouterの内部で行うようにする
	mux := router.NewRouter(todoDB, userID, password, &wg)

	// HTTPサーバーを設定
	srv := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	// シグナルを受け取るためのコンテキストを作成
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// サーバーを別のゴルーチンで起動
	go func() {
		log.Printf("Starting server on %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	// シグナルを待機
	<-ctx.Done()
	log.Println("Shutdown signal received")

	// シャットダウン用のコンテキスト（タイムアウト付き）を作成
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// サーバーのシャットダウンを開始
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	// すべてのリクエストの完了を待つ
	wg.Wait()
	log.Println("Server exited properly")

	return nil
}
