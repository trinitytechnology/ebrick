package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var DefaultServer HttpServer = NewHttpServer()

type HttpServer interface {
	GetRouter() *gin.Engine
	Start() error
}

type httpServer struct {
	opts Options
}

// GetRouter implements HttpServer.
func (h *httpServer) GetRouter() *gin.Engine {
	return h.opts.Router
}

// Start implements HttpServer.
func (h *httpServer) Start() error {
	// Start the Gin server
	h.opts.Logger.Info("Starting HTTP Server")
	if err := h.opts.Router.Run(fmt.Sprintf(":%d", h.opts.Port)); err != nil {
		h.opts.Logger.Fatal("Failed to start Gin server", zap.Error(err))
	}
	return nil
}

func NewHttpServer(opts ...Option) HttpServer {
	return &httpServer{
		opts: newOptions(opts...),
	}
}
