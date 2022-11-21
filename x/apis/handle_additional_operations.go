package apis

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// RunAdditionalOperations implements modules.AdditionalOperationsModule
func (m Module) RunAdditionalOperations() error {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(m.Logger(), gin.Recovery())

	// Register the endpoints
	err := m.registrar(router)
	if err != nil {
		panic(err)
	}

	// Build the HTTP server to be able to shut it down if needed
	httpServer := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", m.cfg.Address, m.cfg.Port),
		Handler:           router,
		ReadHeaderTimeout: time.Minute,
	}

	// Listen for and trap any OS signal to gracefully shutdown and exit
	go m.trapSignal(httpServer)

	// Start the HTTP server
	go m.startServer(httpServer)

	// Block main process (signal capture will call WaitGroup's Done)
	log.Info().Str("module", "apis").Msg("started API server")
	return nil
}

func (m Module) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		log.Debug().Str("module", "apis").Str("path", c.Request.URL.Path).Msg("received request")
	}
}

func (m Module) trapSignal(httpServer *http.Server) {
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)

	// Kill (no param) default send syscall.SIGTERM
	// Kill -2 is syscall.SIGINT
	// Kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Debug().Str("module", "apis").Msg("shutting down API server")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Str("module", "apis").Msg("API server forces to shutdown")
	}

	log.Debug().Str("module", "apis").Msg("API server shutdown")
}

func (m Module) startServer(httpServer *http.Server) {
	err := httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}
