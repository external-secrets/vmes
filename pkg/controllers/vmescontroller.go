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
	"github.com/gusfcarvalho/saferun/pkg/exec"
)

// Reconciler reconciles a vmes.
type Reconciler struct {
	ControllerClass string
	RequeueInterval time.Duration
}

func (r *Reconciler) Reconcile(ctx context.Context) error {
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

	remoteRef := es.Spec.DataFrom[0]

	secretMap, err := secretClient.GetSecretMap(ctx, remoteRef)
	if err != nil {
		return fmt.Errorf("could not get secret: %s, %s %w", remoteRef.Key, es.Name, err)
	}

	if es.Spec.Target.Name == "" {
		es.Spec.Target.Name = "/etc/environment"
	}
	providerData, err := getMapfromFile(es.Spec.Target.Name)
	if err != nil {
		return fmt.Errorf("could not get map from file: %w", err)
	}
	if configdata.PublicKeyFilePath != "" {
		encryptedSecretMap, err := encryptProviderData(secretMap, configdata.PublicKeyFilePath)
		if err != nil {
			return fmt.Errorf("could not get encrypted data map: %w", err)
		}
		providerData = esoutils.MergeByteMap(providerData, encryptedSecretMap)
	} else {
		providerData = esoutils.MergeByteMap(providerData, secretMap)
	}

	err = setMapToFile(es.Spec.Target.Name, providerData)
	if err != nil {
		return fmt.Errorf("could not set map to file: %w", err)
	}

	return nil
}

func getMapfromFile(filepath string) (map[string][]byte, error) {
	providerData := make(map[string][]byte)
	envFile, err := ioutil.ReadFile(filepath)
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

func encryptProviderData(providerData map[string][]byte, PublicKeyFilePath string) (map[string][]byte, error) {
	answer := make(map[string][]byte)
	for k, v := range providerData {
		kd := exec.Encrypt(string(v), PublicKeyFilePath)
		answer["SAFE_RUN_"+k] = []byte(kd)
	}
	return answer, nil
}

func setMapToFile(filepath string, providerData map[string][]byte) error {
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
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
