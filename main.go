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
	"github.com/yaninyzwitty/sqs-go/internal/aws/sqs"
	"github.com/yaninyzwitty/sqs-go/internal/controller"
	"github.com/yaninyzwitty/sqs-go/internal/database"
	"github.com/yaninyzwitty/sqs-go/internal/pkg"
	"github.com/yaninyzwitty/sqs-go/internal/repository"
	"github.com/yaninyzwitty/sqs-go/internal/router"
	"github.com/yaninyzwitty/sqs-go/internal/service"
	"github.com/yaninyzwitty/sqs-go/shared"
)

var (
	queueURL = "https://sqs.eu-north-1.amazonaws.com/651706749096/witty-queue"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 32*time.Second)
	defer cancel()

	var cfg pkg.Config
	if err := shared.LoadConfig(&cfg); err != nil {
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

	sqsClient, err := sqs.LoadSQSClient(ctx, cfg.AWS.Region)
	if err != nil {
		slog.Error("failed to load sqs client", "Error", err)
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

	orderRepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepo)
	orderController := controller.NewOrderController(orderService, sqsClient, &queueURL)
	mux := router.NewRouter(*orderController)

	server := &http.Server{
		Addr:    ":" + fmt.Sprintf("%d", cfg.Server.PORT),
		Handler: mux,
	}

	go shared.StartServer(server)

	slog.Info("server is running at port", "port", cfg.Server.PORT)
	quitCH := make(chan os.Signal, 1)
	signal.Notify(quitCH, os.Interrupt)

	<-quitCH
	shared.ShutdownServer(server)
}
