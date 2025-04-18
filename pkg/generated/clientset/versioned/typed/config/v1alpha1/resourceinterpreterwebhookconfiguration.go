/*
Copyright The Karmada Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	context "context"

	configv1alpha1 "github.com/karmada-io/karmada/pkg/apis/config/v1alpha1"
	scheme "github.com/karmada-io/karmada/pkg/generated/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	gentype "k8s.io/client-go/gentype"
)

// ResourceInterpreterWebhookConfigurationsGetter has a method to return a ResourceInterpreterWebhookConfigurationInterface.
// A group's client should implement this interface.
type ResourceInterpreterWebhookConfigurationsGetter interface {
	ResourceInterpreterWebhookConfigurations() ResourceInterpreterWebhookConfigurationInterface
}

// ResourceInterpreterWebhookConfigurationInterface has methods to work with ResourceInterpreterWebhookConfiguration resources.
type ResourceInterpreterWebhookConfigurationInterface interface {
	Create(ctx context.Context, resourceInterpreterWebhookConfiguration *configv1alpha1.ResourceInterpreterWebhookConfiguration, opts v1.CreateOptions) (*configv1alpha1.ResourceInterpreterWebhookConfiguration, error)
	Update(ctx context.Context, resourceInterpreterWebhookConfiguration *configv1alpha1.ResourceInterpreterWebhookConfiguration, opts v1.UpdateOptions) (*configv1alpha1.ResourceInterpreterWebhookConfiguration, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*configv1alpha1.ResourceInterpreterWebhookConfiguration, error)
	List(ctx context.Context, opts v1.ListOptions) (*configv1alpha1.ResourceInterpreterWebhookConfigurationList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *configv1alpha1.ResourceInterpreterWebhookConfiguration, err error)
	ResourceInterpreterWebhookConfigurationExpansion
}

// resourceInterpreterWebhookConfigurations implements ResourceInterpreterWebhookConfigurationInterface
type resourceInterpreterWebhookConfigurations struct {
	*gentype.ClientWithList[*configv1alpha1.ResourceInterpreterWebhookConfiguration, *configv1alpha1.ResourceInterpreterWebhookConfigurationList]
}

// newResourceInterpreterWebhookConfigurations returns a ResourceInterpreterWebhookConfigurations
func newResourceInterpreterWebhookConfigurations(c *ConfigV1alpha1Client) *resourceInterpreterWebhookConfigurations {
	return &resourceInterpreterWebhookConfigurations{
		gentype.NewClientWithList[*configv1alpha1.ResourceInterpreterWebhookConfiguration, *configv1alpha1.ResourceInterpreterWebhookConfigurationList](
			"resourceinterpreterwebhookconfigurations",
			c.RESTClient(),
			scheme.ParameterCodec,
			"",
			func() *configv1alpha1.ResourceInterpreterWebhookConfiguration {
				return &configv1alpha1.ResourceInterpreterWebhookConfiguration{}
			},
			func() *configv1alpha1.ResourceInterpreterWebhookConfigurationList {
				return &configv1alpha1.ResourceInterpreterWebhookConfigurationList{}
			},
		),
	}
}
