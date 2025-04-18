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
	v1alpha1 "github.com/karmada-io/karmada/pkg/apis/policy/v1alpha1"
	policyv1alpha1 "github.com/karmada-io/karmada/pkg/generated/clientset/versioned/typed/policy/v1alpha1"
	gentype "k8s.io/client-go/gentype"
)

// fakeFederatedResourceQuotas implements FederatedResourceQuotaInterface
type fakeFederatedResourceQuotas struct {
	*gentype.FakeClientWithList[*v1alpha1.FederatedResourceQuota, *v1alpha1.FederatedResourceQuotaList]
	Fake *FakePolicyV1alpha1
}

func newFakeFederatedResourceQuotas(fake *FakePolicyV1alpha1, namespace string) policyv1alpha1.FederatedResourceQuotaInterface {
	return &fakeFederatedResourceQuotas{
		gentype.NewFakeClientWithList[*v1alpha1.FederatedResourceQuota, *v1alpha1.FederatedResourceQuotaList](
			fake.Fake,
			namespace,
			v1alpha1.SchemeGroupVersion.WithResource("federatedresourcequotas"),
			v1alpha1.SchemeGroupVersion.WithKind("FederatedResourceQuota"),
			func() *v1alpha1.FederatedResourceQuota { return &v1alpha1.FederatedResourceQuota{} },
			func() *v1alpha1.FederatedResourceQuotaList { return &v1alpha1.FederatedResourceQuotaList{} },
			func(dst, src *v1alpha1.FederatedResourceQuotaList) { dst.ListMeta = src.ListMeta },
			func(list *v1alpha1.FederatedResourceQuotaList) []*v1alpha1.FederatedResourceQuota {
				return gentype.ToPointerSlice(list.Items)
			},
			func(list *v1alpha1.FederatedResourceQuotaList, items []*v1alpha1.FederatedResourceQuota) {
				list.Items = gentype.FromPointerSlice(items)
			},
		),
		fake,
	}
}
