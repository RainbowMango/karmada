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

package antiaffinity

import (
	"context"

	clusterv1alpha1 "github.com/karmada-io/karmada/pkg/apis/cluster/v1alpha1"
	workv1alpha2 "github.com/karmada-io/karmada/pkg/apis/work/v1alpha2"
	"github.com/karmada-io/karmada/pkg/features"

	"github.com/karmada-io/karmada/pkg/scheduler/cache"
	"github.com/karmada-io/karmada/pkg/scheduler/framework"
)

const (
	// Name is the name of the plugin used in the plugin registry and configurations.
	Name = "AntiAffinity"
)

// AntiAffinity checks if a resourcebinding will violate it's anti-affinity term by being schedule to this cluster.
type AntiAffinity struct{}

var _ framework.FilterPlugin = &AntiAffinity{}

// New instantiates the AntiAffinity plugin.
func New() (framework.Plugin, error) {
	return &AntiAffinity{}, nil
}

// Name returns the plugin name.
func (p *AntiAffinity) Name() string {
	return Name
}

// Filter checks whether scheduling this ResourceBinding to the given cluster
// would violate its anti-affinity constraints against peer ResourceBindings.
func (p *AntiAffinity) Filter(
	_ context.Context,
	bindingSpec *workv1alpha2.ResourceBindingSpec,
	_ *workv1alpha2.ResourceBindingStatus,
	cluster *clusterv1alpha1.Cluster,
	snapshot *cache.Snapshot,
) *framework.Result {
	if !features.FeatureGate.Enabled(features.WorkloadAffinity) {
		return framework.NewResult(framework.Success)
	}

	if bindingSpec.Placement.WorkloadAffinity == nil {
		// WorkloadAffinity is not being used
		return framework.NewResult(framework.Success)
	}

	if snapshot == nil {
		return framework.NewResult(framework.Error, "anti-affinity snapshot is nil")
	}

	antiAffinityGroupLabel := bindingSpec.Placement.WorkloadAffinity.AffinityLabelKey
	antiAffinityLabelValue := bindingSpec.Resource.AffinityGroupLabel[antiAffinityGroupLabel]
	if antiAffinityLabelValue == "" {
		return framework.NewResult(framework.Success)
	}

	thisID := bindingSpec.Resource.Namespace + "/" + bindingSpec.Resource.Name
	peerRBs := snapshot.GetPeerResourceBindings(bindingSpec.Resource.Namespace, bindingSpec.Placement.WorkloadAffinity.AffinityLabelKey, antiAffinityLabelValue)
	for _, peer := range peerRBs {
		if peer == thisID {
			continue
		}
		forbiddenClusters := snapshot.GetClustersForResourceBinding(peer)
		if forbiddenClusters.Has(cluster.Name) {
			return framework.NewResult(framework.Unschedulable, "cluster violates this resource bindings anti-affinity term")
		}
	}

	return framework.NewResult(framework.Success)
}
