// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/karmada-io/karmada/pkg/generated/clientset/versioned/typed/propagationstrategy/v1alpha1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakePropagationstrategyV1alpha1 struct {
	*testing.Fake
}

func (c *FakePropagationstrategyV1alpha1) Overrides(namespace string) v1alpha1.OverrideInterface {
	return &FakeOverrides{c, namespace}
}

func (c *FakePropagationstrategyV1alpha1) PropagationBindings(namespace string) v1alpha1.PropagationBindingInterface {
	return &FakePropagationBindings{c, namespace}
}

func (c *FakePropagationstrategyV1alpha1) PropagationPolicies(namespace string) v1alpha1.PropagationPolicyInterface {
	return &FakePropagationPolicies{c, namespace}
}

func (c *FakePropagationstrategyV1alpha1) PropagationWorks(namespace string) v1alpha1.PropagationWorkInterface {
	return &FakePropagationWorks{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakePropagationstrategyV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
