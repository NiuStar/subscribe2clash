package boot

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"

	"subscribe2clash/app/router"
	"subscribe2clash/internal/global"
)

func initHttpServer() {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	router.RegisterRouter(r)

	srv := &http.Server{
		Addr:    global.Listen,
		Handler: r,
	}

	closed := make(chan struct{})
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit
		log.Println("shutdown server ...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("server shutdown: %v", err)
		}

		close(closed)
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %v\n", err)
	}

	<-closed
}
