package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pro-assistance-dev/sprob/config"
)

type HTTP struct {
	HTTPS        bool
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func NewHTTP(config config.Server) *HTTP {
	return &HTTP{Host: config.Host, Port: config.Port, ReadTimeout: time.Duration(config.ReadTimeout), WriteTimeout: time.Duration(config.WriteTimeout), HTTPS: config.HTTPS}
}

func (i *HTTP) SetFileHeaders(c *gin.Context, fileName string) {
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
}

func (i *HTTP) ListenAndServe(handler http.Handler) {
	srv := &http.Server{
		ReadTimeout:  i.ReadTimeout * time.Second,
		WriteTimeout: i.WriteTimeout * time.Second,
		Handler:      handler,
		Addr:         fmt.Sprintf(":%s", i.Port),
	}

	go func() {
		var err error
		if i.HTTPS {
			err = srv.ListenAndServeTLS("localhost.crt", "localhost.key")
		} else {
			err = srv.ListenAndServe()
		}
		if err != nil {
			log.Fatalln(err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}
}

// func (i *HTTP) CORSMiddleware() gin.HandlerFunc {
// 	// return i.middleware.corsMiddleware()
// }
