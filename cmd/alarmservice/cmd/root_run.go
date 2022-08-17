package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	// "github.com/jmoiron/sqlx"
	// "github.com/pkg/errors"
	// "github.com/yurttasutkan/alarmservice/internal/integration"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/resolver"

	"github.com/yurttasutkan/alarmservice/internal/config"
)

func run(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	tasks := []func() error{
		setLogLevel,
		setSyslog,
		setGRPCResolver,
		printStartMessage,
		// migrateGatewayStats,
		// migrateToClusterKeys,
		// setupIntegration,
		// setupCodec,
		// handleDataDownPayloads,
		// startGatewayPing,
		// setupAPI,
		// setupMonitoring,
	}

	for _, t := range tasks {
		if err := t(); err != nil {
			log.Fatal(err)
		}
	}

	sigChan := make(chan os.Signal)
	exitChan := make(chan struct{})
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	log.WithField("signal", <-sigChan).Info("signal received")
	go func() {
		log.Warning("stopping alarmservice")
		// todo: handle graceful shutdown?
		exitChan <- struct{}{}
	}()
	select {
	case <-exitChan:
	case s := <-sigChan:
		log.WithField("signal", s).Info("signal received, stopping immediately")
	}

	return nil
}

func setLogLevel() error {
	log.SetLevel(log.Level(uint8(config.C.General.LogLevel)))
	return nil
}

func setGRPCResolver() error {
	resolver.SetDefaultScheme(config.C.General.GRPCDefaultResolverScheme)
	return nil
}

func printStartMessage() error {
	log.WithFields(log.Fields{
		"version": version,
		"docs":    "https://www.chirpstack.io/",
	}).Info("starting Alarm Server")
	return nil
}

// func setupIntegration() error {
// 	if err := integration.Setup(config.C); err != nil {
// 		return errors.Wrap(err, "setup integration error")
// 	}

// 	return nil
// }