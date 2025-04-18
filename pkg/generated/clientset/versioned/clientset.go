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

package versioned

import (
	fmt "fmt"
	http "net/http"

	appsv1alpha1 "github.com/karmada-io/karmada/pkg/generated/clientset/versioned/typed/apps/v1alpha1"
	autoscalingv1alpha1 "github.com/karmada-io/karmada/pkg/generated/clientset/versioned/typed/autoscaling/v1alpha1"
	clusterv1alpha1 "github.com/karmada-io/karmada/pkg/generated/clientset/versioned/typed/cluster/v1alpha1"
	configv1alpha1 "github.com/karmada-io/karmada/pkg/generated/clientset/versioned/typed/config/v1alpha1"
	networkingv1alpha1 "github.com/karmada-io/karmada/pkg/generated/clientset/versioned/typed/networking/v1alpha1"
	policyv1alpha1 "github.com/karmada-io/karmada/pkg/generated/clientset/versioned/typed/policy/v1alpha1"
	remedyv1alpha1 "github.com/karmada-io/karmada/pkg/generated/clientset/versioned/typed/remedy/v1alpha1"
	searchv1alpha1 "github.com/karmada-io/karmada/pkg/generated/clientset/versioned/typed/search/v1alpha1"
	workv1alpha1 "github.com/karmada-io/karmada/pkg/generated/clientset/versioned/typed/work/v1alpha1"
	workv1alpha2 "github.com/karmada-io/karmada/pkg/generated/clientset/versioned/typed/work/v1alpha2"
	discovery "k8s.io/client-go/discovery"
	rest "k8s.io/client-go/rest"
	flowcontrol "k8s.io/client-go/util/flowcontrol"
)

type Interface interface {
	Discovery() discovery.DiscoveryInterface
	AppsV1alpha1() appsv1alpha1.AppsV1alpha1Interface
	AutoscalingV1alpha1() autoscalingv1alpha1.AutoscalingV1alpha1Interface
	ClusterV1alpha1() clusterv1alpha1.ClusterV1alpha1Interface
	ConfigV1alpha1() configv1alpha1.ConfigV1alpha1Interface
	NetworkingV1alpha1() networkingv1alpha1.NetworkingV1alpha1Interface
	PolicyV1alpha1() policyv1alpha1.PolicyV1alpha1Interface
	RemedyV1alpha1() remedyv1alpha1.RemedyV1alpha1Interface
	SearchV1alpha1() searchv1alpha1.SearchV1alpha1Interface
	WorkV1alpha1() workv1alpha1.WorkV1alpha1Interface
	WorkV1alpha2() workv1alpha2.WorkV1alpha2Interface
}

// Clientset contains the clients for groups.
type Clientset struct {
	*discovery.DiscoveryClient
	appsV1alpha1        *appsv1alpha1.AppsV1alpha1Client
	autoscalingV1alpha1 *autoscalingv1alpha1.AutoscalingV1alpha1Client
	clusterV1alpha1     *clusterv1alpha1.ClusterV1alpha1Client
	configV1alpha1      *configv1alpha1.ConfigV1alpha1Client
	networkingV1alpha1  *networkingv1alpha1.NetworkingV1alpha1Client
	policyV1alpha1      *policyv1alpha1.PolicyV1alpha1Client
	remedyV1alpha1      *remedyv1alpha1.RemedyV1alpha1Client
	searchV1alpha1      *searchv1alpha1.SearchV1alpha1Client
	workV1alpha1        *workv1alpha1.WorkV1alpha1Client
	workV1alpha2        *workv1alpha2.WorkV1alpha2Client
}

// AppsV1alpha1 retrieves the AppsV1alpha1Client
func (c *Clientset) AppsV1alpha1() appsv1alpha1.AppsV1alpha1Interface {
	return c.appsV1alpha1
}

// AutoscalingV1alpha1 retrieves the AutoscalingV1alpha1Client
func (c *Clientset) AutoscalingV1alpha1() autoscalingv1alpha1.AutoscalingV1alpha1Interface {
	return c.autoscalingV1alpha1
}

// ClusterV1alpha1 retrieves the ClusterV1alpha1Client
func (c *Clientset) ClusterV1alpha1() clusterv1alpha1.ClusterV1alpha1Interface {
	return c.clusterV1alpha1
}

// ConfigV1alpha1 retrieves the ConfigV1alpha1Client
func (c *Clientset) ConfigV1alpha1() configv1alpha1.ConfigV1alpha1Interface {
	return c.configV1alpha1
}

// NetworkingV1alpha1 retrieves the NetworkingV1alpha1Client
func (c *Clientset) NetworkingV1alpha1() networkingv1alpha1.NetworkingV1alpha1Interface {
	return c.networkingV1alpha1
}

// PolicyV1alpha1 retrieves the PolicyV1alpha1Client
func (c *Clientset) PolicyV1alpha1() policyv1alpha1.PolicyV1alpha1Interface {
	return c.policyV1alpha1
}

// RemedyV1alpha1 retrieves the RemedyV1alpha1Client
func (c *Clientset) RemedyV1alpha1() remedyv1alpha1.RemedyV1alpha1Interface {
	return c.remedyV1alpha1
}

// SearchV1alpha1 retrieves the SearchV1alpha1Client
func (c *Clientset) SearchV1alpha1() searchv1alpha1.SearchV1alpha1Interface {
	return c.searchV1alpha1
}

// WorkV1alpha1 retrieves the WorkV1alpha1Client
func (c *Clientset) WorkV1alpha1() workv1alpha1.WorkV1alpha1Interface {
	return c.workV1alpha1
}

// WorkV1alpha2 retrieves the WorkV1alpha2Client
func (c *Clientset) WorkV1alpha2() workv1alpha2.WorkV1alpha2Interface {
	return c.workV1alpha2
}

// Discovery retrieves the DiscoveryClient
func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	if c == nil {
		return nil
	}
	return c.DiscoveryClient
}

// NewForConfig creates a new Clientset for the given config.
// If config's RateLimiter is not set and QPS and Burst are acceptable,
// NewForConfig will generate a rate-limiter in configShallowCopy.
// NewForConfig is equivalent to NewForConfigAndClient(c, httpClient),
// where httpClient was generated with rest.HTTPClientFor(c).
func NewForConfig(c *rest.Config) (*Clientset, error) {
	configShallowCopy := *c

	if configShallowCopy.UserAgent == "" {
		configShallowCopy.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	// share the transport between all clients
	httpClient, err := rest.HTTPClientFor(&configShallowCopy)
	if err != nil {
		return nil, err
	}

	return NewForConfigAndClient(&configShallowCopy, httpClient)
}

// NewForConfigAndClient creates a new Clientset for the given config and http client.
// Note the http client provided takes precedence over the configured transport values.
// If config's RateLimiter is not set and QPS and Burst are acceptable,
// NewForConfigAndClient will generate a rate-limiter in configShallowCopy.
func NewForConfigAndClient(c *rest.Config, httpClient *http.Client) (*Clientset, error) {
	configShallowCopy := *c
	if configShallowCopy.RateLimiter == nil && configShallowCopy.QPS > 0 {
		if configShallowCopy.Burst <= 0 {
			return nil, fmt.Errorf("burst is required to be greater than 0 when RateLimiter is not set and QPS is set to greater than 0")
		}
		configShallowCopy.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(configShallowCopy.QPS, configShallowCopy.Burst)
	}

	var cs Clientset
	var err error
	cs.appsV1alpha1, err = appsv1alpha1.NewForConfigAndClient(&configShallowCopy, httpClient)
	if err != nil {
		return nil, err
	}
	cs.autoscalingV1alpha1, err = autoscalingv1alpha1.NewForConfigAndClient(&configShallowCopy, httpClient)
	if err != nil {
		return nil, err
	}
	cs.clusterV1alpha1, err = clusterv1alpha1.NewForConfigAndClient(&configShallowCopy, httpClient)
	if err != nil {
		return nil, err
	}
	cs.configV1alpha1, err = configv1alpha1.NewForConfigAndClient(&configShallowCopy, httpClient)
	if err != nil {
		return nil, err
	}
	cs.networkingV1alpha1, err = networkingv1alpha1.NewForConfigAndClient(&configShallowCopy, httpClient)
	if err != nil {
		return nil, err
	}
	cs.policyV1alpha1, err = policyv1alpha1.NewForConfigAndClient(&configShallowCopy, httpClient)
	if err != nil {
		return nil, err
	}
	cs.remedyV1alpha1, err = remedyv1alpha1.NewForConfigAndClient(&configShallowCopy, httpClient)
	if err != nil {
		return nil, err
	}
	cs.searchV1alpha1, err = searchv1alpha1.NewForConfigAndClient(&configShallowCopy, httpClient)
	if err != nil {
		return nil, err
	}
	cs.workV1alpha1, err = workv1alpha1.NewForConfigAndClient(&configShallowCopy, httpClient)
	if err != nil {
		return nil, err
	}
	cs.workV1alpha2, err = workv1alpha2.NewForConfigAndClient(&configShallowCopy, httpClient)
	if err != nil {
		return nil, err
	}

	cs.DiscoveryClient, err = discovery.NewDiscoveryClientForConfigAndClient(&configShallowCopy, httpClient)
	if err != nil {
		return nil, err
	}
	return &cs, nil
}

// NewForConfigOrDie creates a new Clientset for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *Clientset {
	cs, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return cs
}

// New creates a new Clientset for the given RESTClient.
func New(c rest.Interface) *Clientset {
	var cs Clientset
	cs.appsV1alpha1 = appsv1alpha1.New(c)
	cs.autoscalingV1alpha1 = autoscalingv1alpha1.New(c)
	cs.clusterV1alpha1 = clusterv1alpha1.New(c)
	cs.configV1alpha1 = configv1alpha1.New(c)
	cs.networkingV1alpha1 = networkingv1alpha1.New(c)
	cs.policyV1alpha1 = policyv1alpha1.New(c)
	cs.remedyV1alpha1 = remedyv1alpha1.New(c)
	cs.searchV1alpha1 = searchv1alpha1.New(c)
	cs.workV1alpha1 = workv1alpha1.New(c)
	cs.workV1alpha2 = workv1alpha2.New(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClient(c)
	return &cs
}
