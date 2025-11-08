/*
Copyright 2022 The Karmada Authors.

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

package options

import (
	"fmt"

	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	apiserverfeatures "k8s.io/apiserver/pkg/features"
	basecompatibility "k8s.io/component-base/compatibility"
)

// Validate validates Options.
func (o *Options) Validate() error {
	var errs []error
	errs = append(errs, o.Etcd.Validate()...)
	errs = append(errs, o.SecureServing.Validate()...)
	errs = append(errs, o.Authentication.Validate()...)
	errs = append(errs, o.Authorization.Validate()...)
	errs = append(errs, o.Audit.Validate()...)
	errs = append(errs, o.Features.Validate()...)
	errs = append(errs, o.CoreAPI.Validate()...)
	errs = append(errs, o.ServerRunOptions.Validate()...)

	// Validate that ListFromCacheSnapshot feature gate is not enabled.
	// ListFromCacheSnapshot is not supported in karmada-search due to consistency issues
	// when serving LIST requests from cache snapshots in a multi-cluster proxy scenario.
	if o.ServerRunOptions.ComponentGlobalsRegistry != nil {
		featureGate := o.ServerRunOptions.ComponentGlobalsRegistry.FeatureGateFor(basecompatibility.DefaultKubeComponent)
		if featureGate != nil && featureGate.Enabled(apiserverfeatures.ListFromCacheSnapshot) {
			errs = append(errs, fmt.Errorf("ListFromCacheSnapshot feature gate is not supported in karmada-search. "+
				"Please ensure --feature-gates does not enable ListFromCacheSnapshot, or use --feature-gates=%s=false",
				apiserverfeatures.ListFromCacheSnapshot))
		}
	}

	return utilerrors.NewAggregate(errs)
}
