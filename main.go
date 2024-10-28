package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/yaninyzwitty/sqs-go/internal/aws/sns"
	"github.com/yaninyzwitty/sqs-go/internal/database"
	"github.com/yaninyzwitty/sqs-go/internal/pkg"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 32*time.Second)
	defer cancel()
	var cfg pkg.Config
	file, err := os.Open("config.yaml")
	if err != nil {
		slog.Error("error opening yaml file", "error", err)
		os.Exit(1)
	}

	defer file.Close()
	if err := cfg.LoadConfig(file); err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	db, err := database.NewDatabaseConnection(ctx, cfg.Database.DATABASE_URL)
	if err != nil {
		slog.Error("Failed to create a database connection", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	err = database.PingDatabase(ctx, db)
	if err != nil {
		slog.Error("Failed to ping database", "error", err)
	}
	os.Exit(1)

	snsClient, err := sns.LoadSnsConfig(ctx, cfg.AWS.Region)
	if err != nil {
		slog.Error("failed to load sns configuration", "error", err)
		os.Exit(1)
	}

	// CREATING A TOPIC IE IN TERRAFORM MIGHT BE A GOOD CHOICE, RATHER THAN THIS APPROACH

	snsTopicArn, err := sns.CreateSnsTopicARN(ctx, "eccomerce-arn", snsClient)
	if err != nil {
		slog.Error("failed to create sns topic arn, topic might already exist", "error", err)

	}
	fmt.Println(snsTopicArn)

}
