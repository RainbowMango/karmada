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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	remedyv1alpha1 "github.com/karmada-io/karmada/pkg/apis/remedy/v1alpha1"
	labels "k8s.io/apimachinery/pkg/labels"
	listers "k8s.io/client-go/listers"
	cache "k8s.io/client-go/tools/cache"
)

// RemedyLister helps list Remedies.
// All objects returned here must be treated as read-only.
type RemedyLister interface {
	// List lists all Remedies in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*remedyv1alpha1.Remedy, err error)
	// Get retrieves the Remedy from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*remedyv1alpha1.Remedy, error)
	RemedyListerExpansion
}

// remedyLister implements the RemedyLister interface.
type remedyLister struct {
	listers.ResourceIndexer[*remedyv1alpha1.Remedy]
}

// NewRemedyLister returns a new RemedyLister.
func NewRemedyLister(indexer cache.Indexer) RemedyLister {
	return &remedyLister{listers.New[*remedyv1alpha1.Remedy](indexer, remedyv1alpha1.Resource("remedy"))}
}
