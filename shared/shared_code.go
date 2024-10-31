package shared

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/yaninyzwitty/sqs-go/internal/pkg"
)

func LoadConfig(cfg *pkg.Config) error {
	file, err := os.Open("config.yaml")
	if err != nil {
		return err
	}
	defer file.Close()
	return cfg.LoadConfig(file)
}

func StartServer(server *http.Server) {
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("failed to start server, server error", "error", err)
		os.Exit(1)
	}
}

func ShutdownServer(server *http.Server) {
	slog.Info("Received termination signal, shutting down server...")
	shutdownCTX, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCTX); err != nil {
		slog.Error("Failed to gracefully shut down server", "error", err)
	}
	slog.Info("Server shutdown successful")
}
