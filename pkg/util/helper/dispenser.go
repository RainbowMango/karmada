/*
Copyright 2025 The Karmada Authors.

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

package helper

import (
	"sort"

	workv1alpha2 "github.com/karmada-io/karmada/pkg/apis/work/v1alpha2"
	"github.com/karmada-io/karmada/pkg/util"
)

// Dispenser aims to divide replicas among clusters by different weights.
type Dispenser struct {
	// Target replicas, should be a positive integer.
	NumReplicas int32
	// Final result.
	Result []workv1alpha2.TargetCluster
}

// NewDispenser will construct a dispenser with target replicas and a prescribed initial result.
func NewDispenser(numReplicas int32, init []workv1alpha2.TargetCluster) *Dispenser {
	cp := make([]workv1alpha2.TargetCluster, len(init))
	copy(cp, init)
	return &Dispenser{NumReplicas: numReplicas, Result: cp}
}

// Done indicates whether finish dispensing.
func (a *Dispenser) Done() bool {
	return a.NumReplicas == 0 && len(a.Result) != 0
}

// TakeByWeight divide replicas by a weight list and merge the result into previous result.
func (a *Dispenser) TakeByWeight(w ClusterWeightInfoList) {
	if a.Done() {
		return
	}
	sum := w.GetWeightSum()
	if sum == 0 {
		return
	}

	sort.Sort(w)

	result := make([]workv1alpha2.TargetCluster, 0, w.Len())
	remain := a.NumReplicas
	for _, info := range w {
		replicas := int32(info.Weight * int64(a.NumReplicas) / sum) // #nosec G115: integer overflow conversion int64 -> int32
		result = append(result, workv1alpha2.TargetCluster{
			Name:     info.ClusterName,
			Replicas: replicas,
		})
		remain -= replicas
	}
	// TODO(Garrybest): take rest replicas by fraction part
	for i := range result {
		if remain == 0 {
			break
		}
		result[i].Replicas++
		remain--
	}

	a.NumReplicas = remain
	a.Result = util.MergeTargetClusters(a.Result, result)
}
