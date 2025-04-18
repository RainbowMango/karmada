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
	v1alpha1 "github.com/karmada-io/karmada/pkg/apis/work/v1alpha1"
	workv1alpha1 "github.com/karmada-io/karmada/pkg/generated/clientset/versioned/typed/work/v1alpha1"
	gentype "k8s.io/client-go/gentype"
)

// fakeWorks implements WorkInterface
type fakeWorks struct {
	*gentype.FakeClientWithList[*v1alpha1.Work, *v1alpha1.WorkList]
	Fake *FakeWorkV1alpha1
}

func newFakeWorks(fake *FakeWorkV1alpha1, namespace string) workv1alpha1.WorkInterface {
	return &fakeWorks{
		gentype.NewFakeClientWithList[*v1alpha1.Work, *v1alpha1.WorkList](
			fake.Fake,
			namespace,
			v1alpha1.SchemeGroupVersion.WithResource("works"),
			v1alpha1.SchemeGroupVersion.WithKind("Work"),
			func() *v1alpha1.Work { return &v1alpha1.Work{} },
			func() *v1alpha1.WorkList { return &v1alpha1.WorkList{} },
			func(dst, src *v1alpha1.WorkList) { dst.ListMeta = src.ListMeta },
			func(list *v1alpha1.WorkList) []*v1alpha1.Work { return gentype.ToPointerSlice(list.Items) },
			func(list *v1alpha1.WorkList, items []*v1alpha1.Work) { list.Items = gentype.FromPointerSlice(items) },
		),
		fake,
	}
}
