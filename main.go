package main

import (
	"context"
	"fmt"

	"time"

	"github.com/external-secrets/vmes/pkg/configdata"
	vmescontroller "github.com/external-secrets/vmes/pkg/controllers"
	"github.com/go-co-op/gocron"
)

func main() {
	fmt.Println("Starting")
	recon := vmescontroller.Reconciler{}
	ctx := context.Background()
	// Reconcile 1 time before scheduler to get refreshInterval
	recon.Reconcile(ctx)
	s := gocron.NewScheduler(time.UTC)
	s.Every(configdata.RefreshInterval.Duration.String()).Do(func(){recon.Reconcile(ctx)})
	s.StartBlocking()
}
