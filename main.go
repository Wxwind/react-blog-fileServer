package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"react-blog-fileServer/middleware"
	"time"
)

var (
	listenAddr string
)

type BindFile struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

func shutdownServer(server *http.Server, logger *log.Logger, quit <-chan os.Signal, done chan<- bool) {
	<-quit
	logger.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}
	close(done)
}

func main() {
	flag.StringVar(&listenAddr, "listen-addr", ":7123", "server listen address")
	flag.Parse()
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)

	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)

	//让quit监听os.Interrupt信号量
	signal.Notify(quit, os.Interrupt)

	router := gin.Default()
	router.Use(middleware.Cors())

	router.StaticFS("/static", http.Dir("assets/"))
	router.POST("/upload/markdown", func(c *gin.Context) {
		var bindFile BindFile

		// Bind file
		if err := c.ShouldBind(&bindFile); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("err: %s", err.Error()))
			return
		}

		// Save uploaded file
		file := bindFile.File
		dst := "assets/markdown/" + filepath.Base(file.Filename)
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}
	})

	server := &http.Server{
		Addr:    listenAddr,
		Handler: router,
	}

	go shutdownServer(server, logger, quit, done)

	logger.Printf("Starting server at %s...\n", listenAddr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %s: %v\n", listenAddr, err)
	}

	<-done
	logger.Println("Server stopped")
}
