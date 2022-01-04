package customclient

import (
	"context"
	"fmt"

	esv1alpha1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1alpha1"
	"github.com/external-secrets/vmes/pkg/configdata"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Client struct {
}

func (c *Client) Get(ctx context.Context, key client.ObjectKey, obj client.Object) error {
	name := key.Name

	switch obj.(type) {
	case *corev1.Secret:
		for _, object := range configdata.ConfigSecret {
			if object.ObjectMeta.Name == name {
				secret, ok := obj.(*corev1.Secret)
				if !ok {
					return fmt.Errorf("secret client did not understand object: %T", obj)
				}
				gvk := secret.GroupVersionKind()
				*secret = object
				secret.SetGroupVersionKind(gvk)
			}
		}
	case *esv1alpha1.ExternalSecret:
		for _, object := range configdata.ConfigES {
			if object.ObjectMeta.Name == name {
				es, ok := obj.(*esv1alpha1.ExternalSecret)
				if !ok {
					return fmt.Errorf("es client did not understand object: %T", obj)
				}
				gvk := es.GroupVersionKind()
				*es = object
				es.SetGroupVersionKind(gvk)
			}
		}
	case *esv1alpha1.SecretStore:
		for _, object := range configdata.ConfigSS {
			if object.ObjectMeta.Name == name {
				ss, ok := obj.(*esv1alpha1.SecretStore)
				if !ok {
					return fmt.Errorf("es client did not understand object: %T", obj)
				}
				gvk := ss.GroupVersionKind()
				*ss = object
				ss.SetGroupVersionKind(gvk)
			}
		}
	default:
		return fmt.Errorf("kind unsuported for customclient")
	}

	return nil
}
func (c *Client) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
	return nil
}
func (c *Client) Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
	return nil
}
func (c *Client) Delete(ctx context.Context, obj client.Object, opts ...client.DeleteOption) error {
	return nil
}
func (c *Client) Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
	return nil
}
func (c *Client) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
	return nil
}
func (c *Client) DeleteAllOf(ctx context.Context, obj client.Object, opts ...client.DeleteAllOfOption) error {
	return nil
}
func (c *Client) Scheme() *runtime.Scheme {
	return nil
}
func (c *Client) RESTMapper() meta.RESTMapper {
	return nil
}
func (c *Client) Status() client.StatusWriter {
	return nil
}
