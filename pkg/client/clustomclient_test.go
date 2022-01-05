package customclient

import (
	"context"
	// "fmt"
	"strings"
	"testing"

	esv1alpha1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1alpha1"
	"github.com/external-secrets/vmes/pkg/configdata"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type customClientTestCase struct {
	keyInput     *client.ObjectKey
	configSecret map[string]corev1.Secret
	configEs     map[string]esv1alpha1.ExternalSecret
	configSS     map[string]esv1alpha1.SecretStore
	Err          error
	expectError  string
	expectedObj  interface{}
}

func makeValidCustomClientTestCase() *customClientTestCase {
	return &customClientTestCase{
		keyInput: &client.ObjectKey{
			Name:      "test",
			Namespace: "default",
		},
		configSecret: make(map[string]corev1.Secret),
		configEs:     make(map[string]esv1alpha1.ExternalSecret),
		configSS:     make(map[string]esv1alpha1.SecretStore),
		Err:          nil,
		expectError:  "",
	}
}

func makeValidCustomClientTestCaseCustom(tweaks ...func(smtc *customClientTestCase)) *customClientTestCase {
	tc := makeValidCustomClientTestCase()
	for _, fn := range tweaks {
		fn(tc)
	}
	return tc
}

func TestCustomClientGetSecret(t *testing.T) {
	testSecret := corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-secret",
			Namespace: "default",
		},
	}
	setSingleConfigSecret := func(ctc *customClientTestCase) {
		ctc.configSecret["secretAc"] = testSecret
		ctc.keyInput = &client.ObjectKey{
			Name:      "test-secret",
			Namespace: "default",
		}
		ctc.expectedObj = testSecret
	}

	setMultiConfigSecret := func(ctc *customClientTestCase) {
		ctc.configSecret["secretAc"] = testSecret
		testSecret2 := testSecret
		testSecret2.ObjectMeta.Name = "test-secret2"
		ctc.configSecret["secretAc2"] = testSecret2
		testSecret3 := testSecret
		testSecret3.ObjectMeta.Name = "test-secret3"
		ctc.configSecret["secretAc3"] = testSecret3
		ctc.keyInput = &client.ObjectKey{
			Name:      "test-secret2",
			Namespace: "default",
		}
		ctc.expectedObj = testSecret2
	}

	setSSsearchSecret := func(ctc *customClientTestCase) {
		ctc.configSS["ss"] = esv1alpha1.SecretStore{
			ObjectMeta: v1.ObjectMeta{
				Name:      "test-ss",
				Namespace: "default",
			},
		}
		ctc.keyInput = &client.ObjectKey{
			Name:      "test-ss",
			Namespace: "default",
		}
		ctc.expectError = "could not find matching resource"
	}

	cases := []*customClientTestCase{
		makeValidCustomClientTestCaseCustom(setSingleConfigSecret),
		makeValidCustomClientTestCaseCustom(setMultiConfigSecret),
		makeValidCustomClientTestCaseCustom(setSSsearchSecret),
	}

	client := Client{}
	for k, v := range cases {
		input := &corev1.Secret{}
		for k, v := range v.configSecret {
			configdata.ConfigSecret[k] = v
		}
		err := client.Get(context.Background(), *v.keyInput, input)
		if !ErrorContains(err, v.expectError) {
			t.Errorf("[%d] unexpected error: %s, expected: '%s'", k, err.Error(), v.expectError)
		}
		if err == nil && input.ObjectMeta.Name != v.expectedObj.(corev1.Secret).ObjectMeta.Name {
			t.Errorf("[%d] unexpected result: expected %v, got %v", k, v.expectedObj.(corev1.Secret), input)
		}
	}
}

func TestCustomClientGetES(t *testing.T) {
	testES := esv1alpha1.ExternalSecret{
		ObjectMeta: v1.ObjectMeta{
			Name: "test-es",
			Namespace: "default",
		},
	}
	setSingleConfigES := func(ctc *customClientTestCase) {
		ctc.configEs["es"] = testES
		ctc.keyInput = &client.ObjectKey{
			Name:      "test-es",
			Namespace: "default",
		}
		ctc.expectedObj = testES
	}

	setMultiConfigES := func(ctc *customClientTestCase) {
		ctc.configEs["es"] = testES
		testES2 := testES
		testES2.ObjectMeta.Name = "test-es2"
		ctc.configEs["es2"] = testES2
		testES3 := testES
		testES3.ObjectMeta.Name = "test-es3"
		ctc.configEs["es3"] = testES3
		ctc.keyInput = &client.ObjectKey{
			Name:      "test-es2",
			Namespace: "default",
		}
		ctc.expectedObj = testES2
	}

	setSSsearchES := func(ctc *customClientTestCase) {
		ctc.configSS["ss"] = esv1alpha1.SecretStore{
			ObjectMeta: v1.ObjectMeta{
				Name:      "test-ss",
				Namespace: "default",
			},
		}
		ctc.keyInput = &client.ObjectKey{
			Name:      "test-ss",
			Namespace: "default",
		}
		ctc.expectError = "could not find matching resource"
	}

	cases := []*customClientTestCase{
		makeValidCustomClientTestCaseCustom(setSingleConfigES),
		makeValidCustomClientTestCaseCustom(setMultiConfigES),
		makeValidCustomClientTestCaseCustom(setSSsearchES),
	}

	client := Client{}
	for k, v := range cases {
		input := &esv1alpha1.ExternalSecret{}
		for k, v := range v.configEs {
			configdata.ConfigES[k] = v
		}
		err := client.Get(context.Background(), *v.keyInput, input)
		if !ErrorContains(err, v.expectError) {
			t.Errorf("[%d] unexpected error: %s, expected: '%s'", k, err.Error(), v.expectError)
		}
		if err == nil && input.ObjectMeta.Name != v.expectedObj.(esv1alpha1.ExternalSecret).ObjectMeta.Name {
			t.Errorf("[%d] unexpected result: expected %v, got %v", k, v.expectedObj.(corev1.Secret), input)
		}
	}
}

func ErrorContains(out error, want string) bool {
	if out == nil {
		return want == ""
	}
	if want == "" {
		return false
	}
	return strings.Contains(out.Error(), want)
}
