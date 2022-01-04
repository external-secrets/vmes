package configdata

import (
	"fmt"
	"os"
	"time"

	// "time"
	"io/ioutil"
	"sync"

	esv1alpha1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1alpha1"
	esmeta "github.com/external-secrets/external-secrets/apis/meta/v1"
	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var ConfigES map[string]esv1alpha1.ExternalSecret
var ConfigSS map[string]esv1alpha1.SecretStore
var ConfigSecret map[string]corev1.Secret
var RefreshInterval v1.Duration
var buildlock sync.RWMutex

type YAMLSecretStore struct {
	Metadata Metadata            `yaml:"metadata"`
	Spec     YAMLSecretStoreSpec `yaml:"spec"`
}

type YAMLSecretStoreSpec struct {
	Provider Provider `yaml:"provider"`
}

type Provider struct {
	// AWS only for now
	AWS AWS `yaml:"aws"`
}

type AWS struct {
	Service string `yaml:"service"`
	Region  string `yaml:"region"`
	Auth    Auth   `yaml:"auth"`
}

type Auth struct {
	SecretRef SecretRef `yaml:"secretRef"`
	Testfield string    `yaml:"testfield"`
}

type SecretRef struct {
	Testfield2               string                   `yaml:"testfield"`
	AccessKeyIDSecretRef     AccessKeyIDSecretRef     `yaml:"accessKeyIDSecretRef"`
	SecretAccessKeySecretRef SecretAccessKeySecretRef `yaml:"secretAccessKeySecretRef"`
}

type AccessKeyIDSecretRef struct {
	Name string `yaml:"name"`
	Key  string `yaml:"key"`
}

type SecretAccessKeySecretRef struct {
	Name string `yaml:"name"`
	Key  string `yaml:"key"`
}

type YAMLExternalSecret struct {
	Metadata Metadata               `yaml:"metadata"`
	Spec     YAMLExternalSecretSpec `yaml:"spec"`
}

type Metadata struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
}

type YAMLExternalSecretSpec struct {
	RefeshInterval string `yaml:"refreshInterval"`
	SecretStoreRef SecretStoreRef `yaml:"secretStoreRef"`
	Target         Target         `yaml:"target"`
	Data           []Data         `yaml:"data"`
	DataFrom       []DataFrom     `yaml:"dataFrom"`
}

type SecretStoreRef struct {
	Name string `yaml:"name"`
	Kind string `yaml:"kind"`
}

type Target struct {
	Name           string `yaml:"name"`
	CreationPolicy string `yaml:"creationPolicy"`
}

type Data struct {
	SecretKey string    `yaml:"secretKey"`
	RemoteRef RemoteRef `yaml:"remoteRef"`
}

type RemoteRef struct {
	Key      string `yaml:"key"`
	Version  string `yaml:"version"`
	Property string `yaml:"property"`
}

type DataFrom struct {
	Key      string `yaml:"key"`
	Version  string `yaml:"version"`
	Property string `yaml:"property"`
}

func YAMLSSToRealSS(y YAMLSecretStore) *esv1alpha1.SecretStore {
	return &esv1alpha1.SecretStore{
		ObjectMeta: v1.ObjectMeta{
			Name:      y.Metadata.Name,
			Namespace: y.Metadata.Namespace,
		},
		Spec: esv1alpha1.SecretStoreSpec{
			Provider: &esv1alpha1.SecretStoreProvider{
				AWS: &esv1alpha1.AWSProvider{
					Region:  y.Spec.Provider.AWS.Region,
					Service: esv1alpha1.AWSServiceType(y.Spec.Provider.AWS.Service),
					Auth: esv1alpha1.AWSAuth{
						SecretRef: &esv1alpha1.AWSAuthSecretRef{
							AccessKeyID: esmeta.SecretKeySelector{
								Name: y.Spec.Provider.AWS.Auth.SecretRef.AccessKeyIDSecretRef.Name,
								Key:  y.Spec.Provider.AWS.Auth.SecretRef.AccessKeyIDSecretRef.Key,
							},
							SecretAccessKey: esmeta.SecretKeySelector{
								Name: y.Spec.Provider.AWS.Auth.SecretRef.SecretAccessKeySecretRef.Name,
								Key:  y.Spec.Provider.AWS.Auth.SecretRef.SecretAccessKeySecretRef.Key,
							},
						},
					},
				},
			},
		},
	}
}

func YAMLEsToRealES(y YAMLExternalSecret) *esv1alpha1.ExternalSecret {
	refreshInterval := v1.Duration{}
	refreshInterval.Duration, _ = time.ParseDuration(y.Spec.RefeshInterval)
	es := &esv1alpha1.ExternalSecret{
		ObjectMeta: v1.ObjectMeta{
			Name:      y.Metadata.Name,
			Namespace: y.Metadata.Namespace,
		},
		Spec: esv1alpha1.ExternalSecretSpec{
			RefreshInterval: &refreshInterval,
			SecretStoreRef: esv1alpha1.SecretStoreRef{
				Name: y.Spec.SecretStoreRef.Name,
				Kind: y.Spec.SecretStoreRef.Kind,
			},
			Target: esv1alpha1.ExternalSecretTarget{
				Name:           y.Spec.Target.Name,
				CreationPolicy: esv1alpha1.ExternalSecretCreationPolicy(y.Spec.Target.CreationPolicy),
			},
		},
	}
	if len(y.Spec.Data) > 0 {
			es.Spec.Data = []esv1alpha1.ExternalSecretData{
			{
				SecretKey: y.Spec.Data[0].SecretKey,
				RemoteRef: esv1alpha1.ExternalSecretDataRemoteRef{
					Key:      y.Spec.Data[0].RemoteRef.Key,
					Version:  y.Spec.Data[0].RemoteRef.Version,
					Property: y.Spec.Data[0].RemoteRef.Property,
				},
			},
		}
	}
	if len(y.Spec.DataFrom) > 0 {
		es.Spec.DataFrom =[]esv1alpha1.ExternalSecretDataRemoteRef{
			{
				Key:      y.Spec.DataFrom[0].Key,
				Version:  y.Spec.DataFrom[0].Version,
				Property: y.Spec.DataFrom[0].Property,
			},
		}
	}
	return es
}


func GetConfigES(s string) (r esv1alpha1.ExternalSecret) {
	buildlock.RLock()
	r = ConfigES[s]
	buildlock.RUnlock()
	return r
}

func GetConfigSS(s string) (r esv1alpha1.SecretStore) {
	buildlock.RLock()
	r = ConfigSS[s]
	buildlock.RUnlock()
	return r
}

func GetConfigSecret(s string) (r corev1.Secret) {
	buildlock.RLock()
	r = ConfigSecret[s]
	buildlock.RUnlock()
	return r
}

func init() {
	duration, _ := time.ParseDuration("10h")
	RefreshInterval = v1.Duration{}
	RefreshInterval.Duration = duration
	ConfigES = make(map[string]esv1alpha1.ExternalSecret)
	ConfigSS = make(map[string]esv1alpha1.SecretStore)
	ConfigSecret = make(map[string]corev1.Secret)
	esyaml := &YAMLExternalSecret{}
	esyamlFile, err := ioutil.ReadFile("pkg/configdata/es.yml")
	if err != nil {
		fmt.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(esyamlFile, esyaml)
	if err != nil {
		fmt.Printf("Unmarshal: %v", err)
	}

	ssyaml := &YAMLSecretStore{}
	ssyamlFile, err := ioutil.ReadFile("pkg/configdata/ss.yml")
	if err != nil {
		fmt.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(ssyamlFile, ssyaml)
	if err != nil {
		fmt.Printf("Unmarshal: %v", err)
	}

	es := YAMLEsToRealES(*esyaml)
	ss := YAMLSSToRealSS(*ssyaml)
	secretAc := corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name:      "awssm-secret",
			Namespace: "default",
		},
		Data: make(map[string][]byte),
	}
	ConfigES["es"] = *es
	ConfigSS["ss"] = *ss
	secretAc.Data["access-key"] = []byte(os.Getenv("AWS_ACCESS_KEY_ID"))
	secretAc.Data["secret-access-key"] = []byte(os.Getenv("AWS_SECRET_ACCESS_KEY"))
	ConfigSecret["secretAc"] = secretAc
}
