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
	v1alpha1 "github.com/karmada-io/karmada/pkg/apis/apps/v1alpha1"
	appsv1alpha1 "github.com/karmada-io/karmada/pkg/generated/clientset/versioned/typed/apps/v1alpha1"
	gentype "k8s.io/client-go/gentype"
)

// fakeWorkloadRebalancers implements WorkloadRebalancerInterface
type fakeWorkloadRebalancers struct {
	*gentype.FakeClientWithList[*v1alpha1.WorkloadRebalancer, *v1alpha1.WorkloadRebalancerList]
	Fake *FakeAppsV1alpha1
}

func newFakeWorkloadRebalancers(fake *FakeAppsV1alpha1) appsv1alpha1.WorkloadRebalancerInterface {
	return &fakeWorkloadRebalancers{
		gentype.NewFakeClientWithList[*v1alpha1.WorkloadRebalancer, *v1alpha1.WorkloadRebalancerList](
			fake.Fake,
			"",
			v1alpha1.SchemeGroupVersion.WithResource("workloadrebalancers"),
			v1alpha1.SchemeGroupVersion.WithKind("WorkloadRebalancer"),
			func() *v1alpha1.WorkloadRebalancer { return &v1alpha1.WorkloadRebalancer{} },
			func() *v1alpha1.WorkloadRebalancerList { return &v1alpha1.WorkloadRebalancerList{} },
			func(dst, src *v1alpha1.WorkloadRebalancerList) { dst.ListMeta = src.ListMeta },
			func(list *v1alpha1.WorkloadRebalancerList) []*v1alpha1.WorkloadRebalancer {
				return gentype.ToPointerSlice(list.Items)
			},
			func(list *v1alpha1.WorkloadRebalancerList, items []*v1alpha1.WorkloadRebalancer) {
				list.Items = gentype.FromPointerSlice(items)
			},
		),
		fake,
	}
}
