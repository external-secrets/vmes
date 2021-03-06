package main

import (
	"context"
	"fmt"

	"flag"
	"time"

	"github.com/external-secrets/vmes/pkg/configdata"
	vmescontroller "github.com/external-secrets/vmes/pkg/controllers"
	"github.com/go-co-op/gocron"
)

func main() {
	flag.StringVar(&configdata.ConfigLocation, "config-path", "/root/.vmes", "Where yaml files should be placed.")
	flag.StringVar(&configdata.PublicKeyFilePath, "public-key-path", "", "Public key file path to use for encryption.")
	flag.Parse()
	configdata.InitConfig()
	fmt.Println("Starting")
	recon := vmescontroller.Reconciler{}
	ctx := context.Background()
	// Reconcile 1 time before scheduler to get refreshInterval
	err := recon.Reconcile(ctx)
	if err != nil {
		fmt.Printf("could not reconcile: %w", err)
	}
	s := gocron.NewScheduler(time.UTC)
	_, err = s.Every(configdata.RefreshInterval.Duration.String()).Do(func() { recon.Reconcile(ctx) })
	if err != nil {
		fmt.Println(err.Error())
	}
	s.StartBlocking()
}
