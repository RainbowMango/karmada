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

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha2

import (
	context "context"
	time "time"

	apisworkv1alpha2 "github.com/karmada-io/karmada/pkg/apis/work/v1alpha2"
	versioned "github.com/karmada-io/karmada/pkg/generated/clientset/versioned"
	internalinterfaces "github.com/karmada-io/karmada/pkg/generated/informers/externalversions/internalinterfaces"
	workv1alpha2 "github.com/karmada-io/karmada/pkg/generated/listers/work/v1alpha2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// ResourceBindingInformer provides access to a shared informer and lister for
// ResourceBindings.
type ResourceBindingInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() workv1alpha2.ResourceBindingLister
}

type resourceBindingInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewResourceBindingInformer constructs a new informer for ResourceBinding type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewResourceBindingInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredResourceBindingInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredResourceBindingInformer constructs a new informer for ResourceBinding type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredResourceBindingInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.WorkV1alpha2().ResourceBindings(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.WorkV1alpha2().ResourceBindings(namespace).Watch(context.TODO(), options)
			},
		},
		&apisworkv1alpha2.ResourceBinding{},
		resyncPeriod,
		indexers,
	)
}

func (f *resourceBindingInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredResourceBindingInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *resourceBindingInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&apisworkv1alpha2.ResourceBinding{}, f.defaultInformer)
}

func (f *resourceBindingInformer) Lister() workv1alpha2.ResourceBindingLister {
	return workv1alpha2.NewResourceBindingLister(f.Informer().GetIndexer())
}
