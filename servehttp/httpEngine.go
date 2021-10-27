package servehttp

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// StartHTTPServer 启动 http 服务
func StartHTTPServer(engine *gin.Engine) {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: engine,
	}

	// 异步运行http服务 (如果服务启动失败 panic 会退出)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// 将调用 os.Exit(1)
			log.Fatalf("listen: %v\n", err)
		}
	}()

	// 监听终止信号
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("[SHUTDOWN] shutdown signal has been received, the service will exit in 3 seconds.")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 优雅终止 http.Server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("[SHUTDOWN] http server shutdown:%v\n", err)
	}
	log.Println("[SHUTDOWN] http server is shutdowning gracefully, new request will be rejected.")

	// waiting for ctx.Done(). timeout of 3 seconds.
	select {
	case <-ctx.Done():
	}

	log.Println("[SHUTDOWN] service exiting")
}
