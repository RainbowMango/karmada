/*
Copyright 2021 The Karmada Authors.

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

package cache

import (
	"sync"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"

	clusterv1alpha1 "github.com/karmada-io/karmada/pkg/apis/cluster/v1alpha1"
	workv1alpha2 "github.com/karmada-io/karmada/pkg/apis/work/v1alpha2"
	"github.com/karmada-io/karmada/pkg/features"
	clusterlister "github.com/karmada-io/karmada/pkg/generated/listers/cluster/v1alpha1"
)

// Cache is an interface for scheduler internal cache.
type Cache interface {
	AddCluster(cluster *clusterv1alpha1.Cluster)
	UpdateCluster(cluster *clusterv1alpha1.Cluster)
	DeleteCluster(cluster *clusterv1alpha1.Cluster)
	// Snapshot returns a snapshot of the current clusters info
	Snapshot() Snapshot
	// Cache should be updated in response to RB changes
	OnResourceBindingAdd(obj interface{})
	OnResourceBindingUpdate(old, cur interface{})
	OnResourceBindingDelete(obj interface{})
}

type AntiKey struct {
	Namespace  string
	LabelKey   string
	GroupValue string
}

func MakeAntiKey(ns, key, value string) AntiKey {
	return AntiKey{ns, key, value}
}

type schedulerCache struct {
	clusterLister  clusterlister.ClusterLister
	mu             sync.RWMutex
	affinityGroups map[AntiKey][]string // antiKey -> []rbID
	clustersByRB   map[string]sets.Set[string]
}

// NewCache instantiates a cache used only by scheduler.
func NewCache(clusterLister clusterlister.ClusterLister) Cache {
	return &schedulerCache{
		clusterLister:  clusterLister,
		affinityGroups: make(map[AntiKey][]string),
		clustersByRB:   make(map[string]sets.Set[string]),
	}
}

// AddCluster does nothing since clusterLister would synchronize automatically
func (c *schedulerCache) AddCluster(_ *clusterv1alpha1.Cluster) {
}

// UpdateCluster does nothing since clusterLister would synchronize automatically
func (c *schedulerCache) UpdateCluster(_ *clusterv1alpha1.Cluster) {
}

// DeleteCluster does nothing since clusterLister would synchronize automatically
func (c *schedulerCache) DeleteCluster(_ *clusterv1alpha1.Cluster) {
}

// Snapshot returns clusters' snapshot.
// **TODO: Needs optimization, only clone when necessary
func (c *schedulerCache) Snapshot() Snapshot {
	out := NewEmptySnapshot()
	clusters, err := c.clusterLister.List(labels.Everything())
	if err != nil {
		klog.Errorf("Failed to list clusters: %v", err)
		return out
	}

	out.clusters = make([]*clusterv1alpha1.Cluster, 0, len(clusters))

	for _, cluster := range clusters {
		out.clusters = append(out.clusters, cluster.DeepCopy())
	}

	// If we have WorkloadAffinity feature enabled, we should index our RBs
	if features.FeatureGate.Enabled(features.WorkloadAffinity) {
		c.mu.RLock()
		defer c.mu.RUnlock()

		out.affinityGroups = make(map[AntiKey][]string, len(c.affinityGroups))
		for k, v := range c.affinityGroups {
			vv := make([]string, len(v))
			copy(vv, v)
			out.affinityGroups[k] = vv
		}

		out.clustersByRB = make(map[string]sets.Set[string], len(c.clustersByRB))
		for rbID, set := range c.clustersByRB {
			s := sets.New[string]()
			for item := range set {
				s.Insert(item)
			}
			out.clustersByRB[rbID] = s
		}
	}

	return out
}

func (c *schedulerCache) OnResourceBindingAdd(obj interface{}) {
	rb := getRBFromObj(obj)
	if rb == nil {
		return
	}
	c.indexRB(rb)
}

func (c *schedulerCache) OnResourceBindingUpdate(oldObj, newObj interface{}) {
	oldRB := getRBFromObj(oldObj)
	newRB := getRBFromObj(newObj)
	if oldRB == nil || newRB == nil {
		return
	}

	c.unindexRB(oldRB)
	c.indexRB(newRB)
}

func (c *schedulerCache) OnResourceBindingDelete(obj interface{}) {
	rb := getRBFromObj(obj)
	if rb == nil {
		return
	}
	c.unindexRB(rb)
}

func (c *schedulerCache) indexRB(rb *workv1alpha2.ResourceBinding) {
	if len(rb.Spec.Clusters) == 0 {
		return
	}

	affinityTerm := rb.Spec.Placement.WorkloadAffinity
	if affinityTerm == nil || affinityTerm.AffinityLabelKey == "" {
		return
	}

	groupValue := rb.Spec.Resource.AffinityGroupLabel[affinityTerm.AffinityLabelKey]
	if groupValue == "" {
		return
	}

	rbID := rb.Namespace + "/" + rb.Name
	clusters := sets.New[string]()
	for _, target := range rb.Spec.Clusters {
		clusters.Insert(target.Name)
	}

	key := MakeAntiKey(rb.Namespace, affinityTerm.AffinityLabelKey, groupValue)

	c.mu.Lock()
	defer c.mu.Unlock()

	c.affinityGroups[key] = append(c.affinityGroups[key], rbID)
	c.clustersByRB[rbID] = clusters
}

func (c *schedulerCache) unindexRB(rb *workv1alpha2.ResourceBinding) {
	affinityTerm := rb.Spec.Placement.WorkloadAffinity
	if affinityTerm == nil || affinityTerm.AffinityLabelKey == "" {
		return
	}

	groupValue := rb.Spec.Resource.AffinityGroupLabel[affinityTerm.AffinityLabelKey]
	if groupValue == "" {
		return
	}

	rbID := rb.Namespace + "/" + rb.Name
	key := MakeAntiKey(rb.Namespace, affinityTerm.AffinityLabelKey, groupValue)

	c.mu.Lock()
	defer c.mu.Unlock()

	if slice, ok := c.affinityGroups[key]; ok {
		filtered := slice[:0]
		for _, id := range slice {
			if id != rbID {
				filtered = append(filtered, id)
			}
		}
		if len(filtered) == 0 {
			delete(c.affinityGroups, key)
		} else {
			c.affinityGroups[key] = filtered
		}
	}

	delete(c.clustersByRB, rbID)
}

func getRBFromObj(obj interface{}) *workv1alpha2.ResourceBinding {
	switch t := obj.(type) {
	case *workv1alpha2.ResourceBinding:
		return t
	case cache.DeletedFinalStateUnknown:
		if rb, ok := t.Obj.(*workv1alpha2.ResourceBinding); ok {
			return rb
		}
	}
	return nil
}
