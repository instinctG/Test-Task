package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Handler struct {
	Router  *gin.Engine
	Service TokenService
	Server  *http.Server
}

func NewHandler(port string, service TokenService) *Handler {
	h := &Handler{Service: service}

	h.Router = gin.Default()
	h.mapRoutes()

	h.Server = &http.Server{
		Addr:    port,
		Handler: h.Router,
	}

	return h
}

func (h *Handler) mapRoutes() {

	h.Router.POST("/token/:id", h.PostRefreshToken)
	h.Router.PUT("refresh/token/:id", h.UpdateRefreshToken)
}

func (h *Handler) Serve(port string) error {
	go func() {
		if err := h.Router.Run(port); err != nil {
			log.Println(err.Error())
		}
	}()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	h.Server.Shutdown(ctx)

	log.Println("shut down gracefully")
	return nil
}
