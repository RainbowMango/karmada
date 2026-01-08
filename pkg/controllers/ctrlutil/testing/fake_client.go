/*
Copyright 2024 The Karmada Authors.

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

// Package testing provides test utilities for controller testing.
// This package should only be used in unit tests.
package testing

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
)

// WithGVKInterceptor returns an interceptor that sets the GVK on objects after Get operations.
// This simulates the behavior of real API server clients, which automatically set GVK.
// The fake client doesn't do this by default (since controller-runtime v0.22.0).
//
// Usage:
//
//	fakeClient := fake.NewClientBuilder().
//	    WithScheme(scheme).
//	    WithInterceptorFuncs(ctrltesting.WithGVKInterceptor(scheme)).
//	    Build()
func WithGVKInterceptor(scheme *runtime.Scheme) interceptor.Funcs {
	return interceptor.Funcs{
		Get: func(ctx context.Context, c client.WithWatch, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
			if err := c.Get(ctx, key, obj, opts...); err != nil {
				return err
			}
			// Set GVK from scheme if it's empty (mimicking real client behavior)
			if obj.GetObjectKind().GroupVersionKind().Empty() {
				gvks, _, _ := scheme.ObjectKinds(obj)
				if len(gvks) > 0 {
					obj.GetObjectKind().SetGroupVersionKind(gvks[0])
				}
			}
			return nil
		},
	}
}
