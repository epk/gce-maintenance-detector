package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"cloud.google.com/go/compute/metadata"
	"go.uber.org/atomic"
	"golang.org/x/sync/errgroup"
)

const (
	maintenanceEvent = "instance/maintenance-event"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	grp, ctx := errgroup.WithContext(ctx)

	// https://cloud.google.com/compute/docs/metadata/getting-live-migration-notice#review_the_outputs
	// starting value: "NONE"
	// can transition to: "MIGRATE_ON_HOST_MAINTENANCE" or "TERMINATE_ON_HOST_MAINTENANCE"
	// will transition back to "NONE" when maintenance is completed
	lastEvent := atomic.NewString("")

	grp.Go(func() error {
		return metadata.SubscribeWithContext(ctx, maintenanceEvent,
			func(ctx context.Context, v string, ok bool) error {
				if !ok {
					logger.Error("failed to get maintenance event", "error", v)
					return nil
				}

				if previousValue := lastEvent.Load(); previousValue != v {
					// Only log when the event value changes
					logger.Info("live maintenance event change", "previous", previousValue, "current", v)

					// carry forward the current value for the next iteration
					lastEvent.Store(v)
				}

				return nil
			},
		)
	})

	if err := grp.Wait(); err != nil && err != context.Canceled {
		logger.Error("errorgroup.Wait returned an error", "error", err)
	}

	logger.Info("exiting")
}
