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

package workloadantiaffinity

import (
	"context"

	"k8s.io/klog/v2"

	clusterv1alpha1 "github.com/karmada-io/karmada/pkg/apis/cluster/v1alpha1"
	workv1alpha2 "github.com/karmada-io/karmada/pkg/apis/work/v1alpha2"
	"github.com/karmada-io/karmada/pkg/scheduler/framework"
)

const (
	// Name is the name of the plugin used in the plugin registry and configurations.
	Name = "WorkloadAntiAffinity"
)

// WorkloadAntiAffinity is a plugin that checks if a cluster matches the workload anti-affinity group constraint.
// It separates workloads in the same anti-affinity group across different clusters.
type WorkloadAntiAffinity struct{}

var _ framework.FilterPlugin = &WorkloadAntiAffinity{}

// New instantiates the workload anti-affinity plugin.
func New() (framework.Plugin, error) {
	return &WorkloadAntiAffinity{}, nil
}

// Name returns the plugin name.
func (p *WorkloadAntiAffinity) Name() string {
	return Name
}

// Filter checks if the cluster matches the workload anti-affinity group constraint.
// If the binding has a workload anti-affinity group specified, it ensures the cluster
// is not selected if it matches the anti-affinity group criteria.
func (p *WorkloadAntiAffinity) Filter(
	_ context.Context,
	bindingSpec *workv1alpha2.ResourceBindingSpec,
	_ *workv1alpha2.ResourceBindingStatus,
	cluster *clusterv1alpha1.Cluster,
) *framework.Result {
	// If no workload anti-affinity groups specified, allow scheduling
	if bindingSpec.WorkloadAffinityGroups == nil || bindingSpec.WorkloadAffinityGroups.AntiAffinityGroup == "" {
		return framework.NewResult(framework.Success)
	}

	antiAffinityGroup := bindingSpec.WorkloadAffinityGroups.AntiAffinityGroup

	// TODO: check if current cluster already have workloads with the same affinity group
	// if yes --> deny
	klog.Infof("checking cluster: %s, with anti-affinity group: %s", cluster.GetName(), antiAffinityGroup)

	// If cluster does not match anti-affinity group label, allow scheduling
	return framework.NewResult(framework.Success)
}
