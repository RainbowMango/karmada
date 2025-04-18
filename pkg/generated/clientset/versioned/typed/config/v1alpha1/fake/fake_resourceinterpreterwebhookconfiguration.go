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

package fake

import (
	v1alpha1 "github.com/karmada-io/karmada/pkg/apis/config/v1alpha1"
	configv1alpha1 "github.com/karmada-io/karmada/pkg/generated/clientset/versioned/typed/config/v1alpha1"
	gentype "k8s.io/client-go/gentype"
)

// fakeResourceInterpreterWebhookConfigurations implements ResourceInterpreterWebhookConfigurationInterface
type fakeResourceInterpreterWebhookConfigurations struct {
	*gentype.FakeClientWithList[*v1alpha1.ResourceInterpreterWebhookConfiguration, *v1alpha1.ResourceInterpreterWebhookConfigurationList]
	Fake *FakeConfigV1alpha1
}

func newFakeResourceInterpreterWebhookConfigurations(fake *FakeConfigV1alpha1) configv1alpha1.ResourceInterpreterWebhookConfigurationInterface {
	return &fakeResourceInterpreterWebhookConfigurations{
		gentype.NewFakeClientWithList[*v1alpha1.ResourceInterpreterWebhookConfiguration, *v1alpha1.ResourceInterpreterWebhookConfigurationList](
			fake.Fake,
			"",
			v1alpha1.SchemeGroupVersion.WithResource("resourceinterpreterwebhookconfigurations"),
			v1alpha1.SchemeGroupVersion.WithKind("ResourceInterpreterWebhookConfiguration"),
			func() *v1alpha1.ResourceInterpreterWebhookConfiguration {
				return &v1alpha1.ResourceInterpreterWebhookConfiguration{}
			},
			func() *v1alpha1.ResourceInterpreterWebhookConfigurationList {
				return &v1alpha1.ResourceInterpreterWebhookConfigurationList{}
			},
			func(dst, src *v1alpha1.ResourceInterpreterWebhookConfigurationList) { dst.ListMeta = src.ListMeta },
			func(list *v1alpha1.ResourceInterpreterWebhookConfigurationList) []*v1alpha1.ResourceInterpreterWebhookConfiguration {
				return gentype.ToPointerSlice(list.Items)
			},
			func(list *v1alpha1.ResourceInterpreterWebhookConfigurationList, items []*v1alpha1.ResourceInterpreterWebhookConfiguration) {
				list.Items = gentype.FromPointerSlice(items)
			},
		),
		fake,
	}
}
