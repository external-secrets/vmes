package vmescontroller

import (
	"context"
	"fmt"

	"io/ioutil"
	"os"
	"strings"
	"time"

	awsprovider "github.com/external-secrets/external-secrets/pkg/provider/aws"
	esoutils "github.com/external-secrets/external-secrets/pkg/utils"
	customclient "github.com/external-secrets/vmes/pkg/client"
	"github.com/external-secrets/vmes/pkg/configdata"
)

// Reconciler reconciles a vmes.
type Reconciler struct {
	ControllerClass string
	RequeueInterval time.Duration
}

func (r *Reconciler) Reconcile(ctx context.Context) error {
	fmt.Println("Reconciling")
	es := configdata.GetConfigES("es")
	if es.Spec.RefreshInterval != nil {
		configdata.RefreshInterval = *es.Spec.RefreshInterval
	}
	store := configdata.GetConfigSS("ss")
	storeProvider := &awsprovider.Provider{}

	client := &customclient.Client{}

	secretClient, err := storeProvider.NewClient(ctx, &store, client, "")
	if err != nil {
		return fmt.Errorf("could not create Client: %w", err)
	}

	remoteRef := es.Spec.Data[0].RemoteRef

	secretMap, err := secretClient.GetSecretMap(ctx, remoteRef)
	if err != nil {
		return fmt.Errorf("could not get secret: %s, %s %w", remoteRef.Key, es.Name, err)
	}

	providerData, err := getMapfromFile()
	if err != nil {
		return fmt.Errorf("could not get map from file: %w", err)
	}

	providerData = esoutils.MergeByteMap(providerData, secretMap)

	err = setMapToFile(providerData)
	if err != nil {
		return fmt.Errorf("could not set map to file: %w", err)
	}

	return nil
}

func getMapfromFile() (map[string][]byte, error) {
	providerData := make(map[string][]byte)
	envFile, err := ioutil.ReadFile("/etc/environment")
	if err != nil {
		return nil, fmt.Errorf("yamlFile.Get err  %w ", err)
	}
	lines := strings.Split(string(envFile), "\n")

	for _, line := range lines {
		if len(line) > 0 && string([]rune(line[0:1])) != "#" {
			splits := strings.Split(line, "=")
			providerData[splits[0]] = []byte(splits[1])
		}
	}
	return providerData, err
}

func setMapToFile(providerData map[string][]byte) error {
	f, err := os.OpenFile("/etc/environment", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	err = f.Truncate(0)
	if err != nil {
		return fmt.Errorf("could not truncate the file: %w", err)
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("could not set file pointer to the beggining: %w", err)
	}
	for k, v := range providerData {
		_, err := f.Write([]byte(fmt.Sprintf("%s=%s\n", k, string(v))))
		if err != nil {
			return fmt.Errorf("could not write to file: %w", err)
		}
	}
	f.Sync()
	return nil
}
