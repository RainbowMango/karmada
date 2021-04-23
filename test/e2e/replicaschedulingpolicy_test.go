package e2e

import (
	"context"
	"fmt"
	"strings"

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
		resourceTemplate := helper.NewDeployment(testNamespace, rand.String(RandomStrLength))
		// we will select two member clusters for this test.
		var selectedClusters []string
		selectedClustersLen := 2

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
				ClusterNames: []string{"member1", "member2"},
			},
		}
		createdClusterPropagationPolicy := helper.NewClusterPropagationPolicy(rand.String(RandomStrLength), selector, placement)
		createdReplicaSchedulingPolicy := &policyv1alpha1.ReplicaSchedulingPolicy{}

		// Deploy ClusterPropagationPolicy
		ginkgo.BeforeEach(func() {
			// We don't care total number of clusters we have, we only select two of them.
			ginkgo.By(fmt.Sprintf("Selecting clusters"), func() {
				selectedClusters = make([]string, 0, 2)
				for i := 0; i < 2 && i < len(clusters); i++ {
					selectedClusters = append(selectedClusters, clusters[i].Name)
				}
				gomega.Expect(selectedClusters).Should(gomega.HaveLen(selectedClustersLen))
				klog.Infof("Selected clusters: %s", strings.Join(selectedClusters, ","))
			})

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
