// Code generated by main. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"time"

	v1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/internalpackaging/v1alpha1"
	scheme "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// InternalPackagesGetter has a method to return a InternalPackageInterface.
// A group's client should implement this interface.
type InternalPackagesGetter interface {
	InternalPackages(namespace string) InternalPackageInterface
}

// InternalPackageInterface has methods to work with InternalPackage resources.
type InternalPackageInterface interface {
	Create(ctx context.Context, internalPackage *v1alpha1.InternalPackage, opts v1.CreateOptions) (*v1alpha1.InternalPackage, error)
	Update(ctx context.Context, internalPackage *v1alpha1.InternalPackage, opts v1.UpdateOptions) (*v1alpha1.InternalPackage, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.InternalPackage, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.InternalPackageList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.InternalPackage, err error)
	InternalPackageExpansion
}

// internalPackages implements InternalPackageInterface
type internalPackages struct {
	client rest.Interface
	ns     string
}

// newInternalPackages returns a InternalPackages
func newInternalPackages(c *InternalV1alpha1Client, namespace string) *internalPackages {
	return &internalPackages{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the internalPackage, and returns the corresponding internalPackage object, and an error if there is any.
func (c *internalPackages) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.InternalPackage, err error) {
	result = &v1alpha1.InternalPackage{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("internalpackages").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of InternalPackages that match those selectors.
func (c *internalPackages) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.InternalPackageList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.InternalPackageList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("internalpackages").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested internalPackages.
func (c *internalPackages) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("internalpackages").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a internalPackage and creates it.  Returns the server's representation of the internalPackage, and an error, if there is any.
func (c *internalPackages) Create(ctx context.Context, internalPackage *v1alpha1.InternalPackage, opts v1.CreateOptions) (result *v1alpha1.InternalPackage, err error) {
	result = &v1alpha1.InternalPackage{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("internalpackages").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(internalPackage).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a internalPackage and updates it. Returns the server's representation of the internalPackage, and an error, if there is any.
func (c *internalPackages) Update(ctx context.Context, internalPackage *v1alpha1.InternalPackage, opts v1.UpdateOptions) (result *v1alpha1.InternalPackage, err error) {
	result = &v1alpha1.InternalPackage{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("internalpackages").
		Name(internalPackage.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(internalPackage).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the internalPackage and deletes it. Returns an error if one occurs.
func (c *internalPackages) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("internalpackages").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *internalPackages) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("internalpackages").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched internalPackage.
func (c *internalPackages) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.InternalPackage, err error) {
	result = &v1alpha1.InternalPackage{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("internalpackages").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}