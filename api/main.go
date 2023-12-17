package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"payhere/api/controller"
	"payhere/config"
	"payhere/repository"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := initGin()
	sig := initSig(r)
	conf := getConfig()

	db := repository.InitDB(conf, "")

	if err := controller.InitHandler(conf, r, db, sig); err != nil {
		os.Exit(1)
	}

	startServer(conf, r, sig)
}

func getConfig() *config.ViperConfig {
	payhere := config.Payhere
	return payhere
}

func initGin() *gin.Engine {
	r := gin.Default()
	r.GET("/healthCheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "alive",
		})
	})
	return r
}

func initSig(r *gin.Engine) chan os.Signal {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go func() {
		<-sc
		signal.Stop(sc)
		close(sc)
	}()

	return sc
}

func startServer(conf *config.ViperConfig, r *gin.Engine, quit <-chan os.Signal) {
	server := fmt.Sprintf("0.0.0.0:%d", conf.GetInt("port"))

	srv := &http.Server{
		Addr:    server,
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("gin start error : ", "err : ", err)
		}
	}()

	<-quit
	time.Sleep(5 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server shutdown : ", err)
	}
}
