package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/yaninyzwitty/sqs-go/internal/aws/sns"
	"github.com/yaninyzwitty/sqs-go/internal/controller"
	"github.com/yaninyzwitty/sqs-go/internal/database"
	"github.com/yaninyzwitty/sqs-go/internal/pkg"
	"github.com/yaninyzwitty/sqs-go/internal/repository"
	"github.com/yaninyzwitty/sqs-go/internal/router"
	"github.com/yaninyzwitty/sqs-go/internal/service"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 32*time.Second)
	defer cancel()

	var cfg pkg.Config
	if err := loadConfig(&cfg); err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	db, err := database.NewDatabaseConnection(ctx, cfg.Database.DATABASE_URL)
	if err != nil {
		slog.Error("Failed to create a database connection", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := database.PingDatabase(ctx, db); err != nil {
		slog.Error("Failed to ping database", "error", err)
		os.Exit(1)
	}

	snsClient, err := sns.LoadSnsConfig(ctx, cfg.AWS.Region)
	if err != nil {
		slog.Error("failed to load sns configuration", "error", err)
		os.Exit(1)
	}

	snsTopicArn, err := sns.CreateSnsTopicARN(ctx, "eccomerce-arn", snsClient)
	if err != nil {
		slog.Error("failed to create sns topic arn, topic might already exist", "error", err)
	}
	fmt.Println(snsTopicArn)

	orderRepo := repository.NewOrderRepository(snsClient, &snsTopicArn, db)
	orderService := service.NewOrderService(orderRepo)
	orderController := controller.NewOrderController(orderService)
	mux := router.NewRouter(*orderController)

	server := &http.Server{
		Addr:    ":" + fmt.Sprintf("%d", cfg.Server.PORT),
		Handler: mux,
	}

	go startServer(server)

	slog.Info("server is running at port", "port", cfg.Server.PORT)
	quitCH := make(chan os.Signal, 1)
	signal.Notify(quitCH, os.Interrupt)

	<-quitCH
	shutdownServer(server)
}

func loadConfig(cfg *pkg.Config) error {
	file, err := os.Open("config.yaml")
	if err != nil {
		return err
	}
	defer file.Close()
	return cfg.LoadConfig(file)
}

func startServer(server *http.Server) {
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("failed to start server, server error", "error", err)
		os.Exit(1)
	}
}

func shutdownServer(server *http.Server) {
	slog.Info("Received termination signal, shutting down server...")
	shutdownCTX, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCTX); err != nil {
		slog.Error("Failed to gracefully shut down server", "error", err)
	}
	slog.Info("Server shutdown successful")
}
