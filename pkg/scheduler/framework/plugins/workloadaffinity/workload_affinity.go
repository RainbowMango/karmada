/*
Copyright 2026 The Karmada Authors.

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

package workloadaffinity

import (
	"context"

	"k8s.io/klog/v2"

	clusterv1alpha1 "github.com/karmada-io/karmada/pkg/apis/cluster/v1alpha1"
	workv1alpha2 "github.com/karmada-io/karmada/pkg/apis/work/v1alpha2"
	"github.com/karmada-io/karmada/pkg/scheduler/framework"
)

const (
	// Name is the name of the plugin used in the plugin registry and configurations.
	Name = "WorkloadAffinity"
)

// WorkloadAffinity is a plugin that checks if a cluster matches the workload affinity group constraint.
// It co-locates workloads in the same affinity group across clusters.
type WorkloadAffinity struct{}

var _ framework.FilterPlugin = &WorkloadAffinity{}

// New instantiates the workload affinity plugin.
func New() (framework.Plugin, error) {
	return &WorkloadAffinity{}, nil
}

// Name returns the plugin name.
func (p *WorkloadAffinity) Name() string {
	return Name
}

// Filter checks if the cluster matches the workload affinity group constraint.
// If the binding has a workload affinity group specified, it ensures the cluster
// is selected based on the affinity group criteria.
func (p *WorkloadAffinity) Filter(
	_ context.Context,
	bindingSpec *workv1alpha2.ResourceBindingSpec,
	_ *workv1alpha2.ResourceBindingStatus,
	cluster *clusterv1alpha1.Cluster,
) *framework.Result {
	// If no workload affinity groups specified, allow scheduling
	if bindingSpec.WorkloadAffinityGroups == nil || bindingSpec.WorkloadAffinityGroups.AffinityGroup == "" {
		return framework.NewResult(framework.Success)
	}

	affinityGroup := bindingSpec.WorkloadAffinityGroups.AffinityGroup

	// TODO: check if current cluster already have workloads with the same affinity group
	// if yes --> allow scheduling
	klog.Infof("checking cluster: %s, with affinity group: %s", cluster.GetName(), affinityGroup)

	// If current cluster doesn't have workloads with the same affinity group, deny
	return framework.NewResult(framework.Unschedulable, "cluster does not match workload affinity group constraint")
}
