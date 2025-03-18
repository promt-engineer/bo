package http

import (
	"backoffice/docs"
	"backoffice/internal/transport/http/middlewares"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	wg     *sync.WaitGroup
	ctx    context.Context
	server *http.Server
	router *gin.Engine
}

// @SecurityDefinitions.apikey X-Authenticate
// @in header
// @name X-Authenticate
func New(ctx context.Context, wg *sync.WaitGroup, cfg *Config, handlers []Handler) *Server {
	docs.SwaggerInfo.Title = "Backoffice API"
	docs.SwaggerInfo.Description = "Backoffice API server."
	docs.SwaggerInfo.Version = "2.0"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	docs.SwaggerInfo.Host = cfg.Domain

	s := &Server{
		wg: wg, ctx: ctx,
		server: &http.Server{
			Addr:              fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Handler:           nil,
			ReadHeaderTimeout: 10 * time.Second,
			ReadTimeout:       cfg.ReadTimeout,
			WriteTimeout:      cfg.WriteTimeout,
			IdleTimeout:       30 * time.Second,
		},
		router: gin.New(),
	}

	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	s.router.Use(middlewares.CORSMiddleware())

	api := s.router.Group("")
	s.registerHandlers(api, handlers...)

	return s
}

func (s *Server) registerHandlers(api *gin.RouterGroup, handlers ...Handler) {
	for _, h := range handlers {
		h.Register(api)
	}

	s.server.Handler = s.router
}

func (s *Server) Run() {
	s.wg.Add(1)
	zap.S().Infof("server listining: %s", s.server.Addr)

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		zap.S().Error(err.Error())
	}
}

func (s *Server) Shutdown() error {
	zap.S().Info("Shutdown server...")
	zap.S().Info("Stopping http server...")

	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
	defer func() {
		cancel()
		s.wg.Done()
	}()

	if err := s.server.Shutdown(ctx); err != nil {
		zap.S().Fatal("Server forced to shutdown:", zap.Error(err))

		return err
	}

	zap.S().Info("Server successfully stopped.")

	return nil
}
