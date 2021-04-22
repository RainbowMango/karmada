package e2e

import (
	"fmt"

	policyv1alpha1 "github.com/karmada-io/karmada/pkg/apis/policy/v1alpha1"
	"github.com/karmada-io/karmada/test/helper"
	"github.com/onsi/ginkgo"
	"k8s.io/apimachinery/pkg/util/rand"
)

var _ = ginkgo.Describe("[ReplicaScheduling] replica scheduling testing", func() {

	// The replicas specified in resource template will be discarded when there is a RSP.
	ginkgo.Context("total replicas should follow the policy", func() {
		resourceTemplate := helper.NewDeployment(testNamespace, rand.String(RandomStrLength))
		// createdClusterPropagationPolicy := &policyv1alpha1.ClusterPropagationPolicy{}
		createdReplicaSchedulingPolicy := &policyv1alpha1.ReplicaSchedulingPolicy{}

		// Deploy ClusterPropagationPolicy
		ginkgo.BeforeEach(func() {

		})

		// Deploy ReplicaSchedulingPolicy
		ginkgo.BeforeEach(func() {

		})

		// Cleanup ReplicaSchedulingPolicy
		ginkgo.AfterEach(func() {

		})

		// Cleanup ClusterPropagationPolicy
		ginkgo.AfterEach(func() {

		})

		ginkgo.By(fmt.Sprintf("Creating deployment(%s/%s)", resourceTemplate.Namespace, resourceTemplate.Name), func() {

		})

		ginkgo.By(fmt.Sprintf("Checking total replicas should be euqal to %d", createdReplicaSchedulingPolicy.Spec.TotalReplicas), func() {

		})

		ginkgo.By(fmt.Sprintf("Deleting deployment(%s/%s)", resourceTemplate.Namespace, resourceTemplate.Name), func() {

		})
	})

	ginkgo.Context("replicas should be allocated via weight list", func() {

	})
})
