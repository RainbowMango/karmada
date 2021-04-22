package e2e

import (
	"context"
	"fmt"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/klog/v2"

	policyv1alpha1 "github.com/karmada-io/karmada/pkg/apis/policy/v1alpha1"
	"github.com/karmada-io/karmada/test/helper"
)

var _ = ginkgo.Describe("[ReplicaScheduling] replica scheduling testing", func() {

	// The replicas specified in resource template will be discarded when there is a RSP.
	ginkgo.Context("total replicas should follow the policy", func() {
		ginkgo.By(fmt.Sprintf("testing clusters"), func() {
			clusterLen := len(clusters)
			if clusterLen < MinimumCluster {
				klog.Errorf("Needs at least %d member clusters to run, but got: %d", MinimumCluster, len(clusters))
			}
		})
		clusterLen := len(clusters)
		if clusterLen < MinimumCluster {
			klog.Errorf("Needs at least %d member clusters to run, but got: %d", MinimumCluster, len(clusters))
		}

		resourceTemplate := helper.NewDeployment(testNamespace, rand.String(RandomStrLength))
		selector := []policyv1alpha1.ResourceSelector{
			{
				APIVersion: resourceTemplate.APIVersion,
				Kind:       resourceTemplate.Kind,
				Namespace:  resourceTemplate.Namespace,
				Name:       resourceTemplate.Name,
			},
		}
		placement := policyv1alpha1.Placement{
			ClusterAffinity: &policyv1alpha1.ClusterAffinity{
				// ClusterNames: []string{clusters[0].ClusterName, clusters[1].ClusterName},
				ClusterNames: []string{},
			},
		}
		createdClusterPropagationPolicy := helper.NewClusterPropagationPolicy(rand.String(RandomStrLength), selector, placement)
		createdReplicaSchedulingPolicy := &policyv1alpha1.ReplicaSchedulingPolicy{}

		// Deploy ClusterPropagationPolicy
		ginkgo.BeforeEach(func() {
			ginkgo.By(fmt.Sprintf("Creating ClusterPropagationPolicy(%s)", createdClusterPropagationPolicy.Name), func() {
				err := controlPlaneClient.Create(context.TODO(), createdClusterPropagationPolicy)
				gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			})
		})

		// Deploy ReplicaSchedulingPolicy
		ginkgo.BeforeEach(func() {

		})

		// Cleanup ReplicaSchedulingPolicy
		ginkgo.AfterEach(func() {
		})

		// Cleanup ClusterPropagationPolicy
		ginkgo.AfterEach(func() {
			ginkgo.By(fmt.Sprintf("Deleting ClusterPropagationPolicy(%s)", createdClusterPropagationPolicy.Name), func() {
				err := controlPlaneClient.Delete(context.TODO(), createdClusterPropagationPolicy)
				gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			})
		})

		ginkgo.It("total replicas should follow the policy", func() {
			ginkgo.By(fmt.Sprintf("Creating deployment(%s/%s)", resourceTemplate.Namespace, resourceTemplate.Name), func() {
				err := controlPlaneClient.Create(context.TODO(), resourceTemplate)
				gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

				clusterLen := len(clusters)
				if clusterLen < MinimumCluster {
					klog.Errorf("Needs at least %d member clusters to run, but got: %d", MinimumCluster, len(clusters))
				}
			})

			ginkgo.By(fmt.Sprintf("Checking total replicas should be euqal to %d", createdReplicaSchedulingPolicy.Spec.TotalReplicas), func() {

			})

			ginkgo.By(fmt.Sprintf("Deleting deployment(%s/%s)", resourceTemplate.Namespace, resourceTemplate.Name), func() {
				err := controlPlaneClient.Delete(context.TODO(), resourceTemplate)
				gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			})
		})

		ginkgo.It("replicas should be allocated via weight list", func() {

		})
	})
})
