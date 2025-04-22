/*
Copyright 2023 The Kubernetes Authors.

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

package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	"k8s.io/utils/net"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"sigs.k8s.io/network-policy-api/apis/v1alpha1"
	"sigs.k8s.io/network-policy-api/conformance/utils/kubernetes"
	"sigs.k8s.io/network-policy-api/conformance/utils/suite"
)

func init() {
	ConformanceTests = append(ConformanceTests,
		AdminNetworkPolicyEgressNamedPort,
		AdminNetworkPolicyEgressNodePeers,
		AdminNetworkPolicyEgressInlineCIDRPeers,
	)
}

var AdminNetworkPolicyEgressNamedPort = suite.ConformanceTest{
	ShortName:   "AdminNetworkPolicyEgressNamedPort",
	Description: "Tests support for egress traffic on a named port using admin network policy API based on a server and client model",
	Features: []suite.SupportedFeature{
		suite.SupportAdminNetworkPolicy,
		suite.SupportAdminNetworkPolicyNamedPorts,
	},
	Manifests: []string{"base/admin_network_policy/standard-egress-tcp-rules.yaml"},
	Test: func(t *testing.T, s *suite.ConformanceTestSuite) {

		t.Run("Should support an 'allow-egress' policy for named port", func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), s.TimeoutConfig.GetTimeout)
			defer cancel()
			// This test uses `egress-tcp` ANP
			// cedric-diggory-1 is our server pod in hufflepuff namespace
			serverPod := &v1.Pod{}
			err := s.Client.Get(ctx, client.ObjectKey{
				Namespace: "network-policy-conformance-hufflepuff",
				Name:      "cedric-diggory-1",
			}, serverPod)
			require.NoErrorf(t, err, "unable to fetch the server pod")
			anp := &v1alpha1.AdminNetworkPolicy{}
			err = s.Client.Get(ctx, client.ObjectKey{
				Name: "egress-tcp",
			}, anp)
			require.NoErrorf(t, err, "unable to fetch the admin network policy")
			mutate := anp.DeepCopy()
			namedPortRule := mutate.Spec.Egress[5]
			webPort := "web"
			// replace the tcp port 8080 rule as named port rule which translate to tcp port 80 instead
			namedPortRule.Ports = &[]v1alpha1.AdminNetworkPolicyPort{
				{
					NamedPort: &webPort,
				},
			}
			mutate.Spec.Egress[5] = namedPortRule
			err = s.Client.Patch(ctx, mutate, client.MergeFrom(anp))
			require.NoErrorf(t, err, "unable to patch the admin network policy")
			// harry-potter-0 is our client pod in gryffindor namespace
			// ensure egress is ALLOWED to hufflepuff from gryffindor at the web port, which is defined as TCP at port 80 in pod spec
			// egressRule at index5 should take effect
			success := kubernetes.PokeServer(t, s.ClientSet, &s.KubeConfig, "network-policy-conformance-gryffindor", "harry-potter-0", "tcp",
				serverPod.Status.PodIP, int32(80), s.TimeoutConfig.RequestTimeout, true)
			assert.True(t, success)
			// harry-potter-1 is our client pod in gryffindor namespace
			// ensure egress is DENIED to hufflepuff from gryffindor for rest of the traffic; egressRule at index6 should take effect
			success = kubernetes.PokeServer(t, s.ClientSet, &s.KubeConfig, "network-policy-conformance-gryffindor", "harry-potter-1", "tcp",
				serverPod.Status.PodIP, int32(8080), s.TimeoutConfig.RequestTimeout, false)
			assert.True(t, success)
		})

	},
}

var AdminNetworkPolicyEgressNodePeers = suite.ConformanceTest{
	ShortName:   "AdminNetworkPolicyEgressNodePeers",
	Description: "Tests support for egress traffic to node peers using admin network policy API based on a server and client model",
	Features: []suite.SupportedFeature{
		suite.SupportAdminNetworkPolicy,
		suite.SupportAdminNetworkPolicyEgressNodePeers,
	},
	Manifests: []string{"base/admin_network_policy/experimental-egress-selector-rules.yaml"},
	Test: func(t *testing.T, s *suite.ConformanceTestSuite) {
		ctx, cancel := context.WithTimeout(context.Background(), s.TimeoutConfig.GetTimeout)
		defer cancel()
		// This test uses `node-and-cidr-as-peers-example` ANP
		// centaur-1 is our server host-networked pod in forbidden-forrest namespace
		serverPod := &v1.Pod{}
		err := s.Client.Get(ctx, client.ObjectKey{
			Namespace: "network-policy-conformance-forbidden-forrest",
			Name:      "centaur-1",
		}, serverPod)
		require.NoErrorf(t, err, "unable to fetch the server pod")
		t.Run("Should support an 'allow-egress' rule policy for egress-node-peer", func(t *testing.T) {
			// harry-potter-0 is our client pod in gryffindor namespace
			// ensure egress is ALLOWED to forbidden-forrest from gryffindor at the 36363 TCP port
			// egressRule at index0 should take effect
			success := kubernetes.PokeServer(t, s.ClientSet, &s.KubeConfig, "network-policy-conformance-gryffindor", "harry-potter-0", "tcp",
				serverPod.Status.PodIP, int32(36363), s.TimeoutConfig.RequestTimeout, true)
			assert.True(t, success)
		})
		t.Run("Should support a 'pass-egress' rule policy for egress-node-peer", func(t *testing.T) {
			// harry-potter-0 is our client pod in gryffindor namespace
			// ensure egress is PASSED to forbidden-forrest from gryffindor at the 34345 UDP port
			// egressRule at index1 should take effect
			success := kubernetes.PokeServer(t, s.ClientSet, &s.KubeConfig, "network-policy-conformance-gryffindor", "harry-potter-1", "udp",
				serverPod.Status.PodIP, int32(34345), s.TimeoutConfig.RequestTimeout, true) // Pass rule at index2 takes effect
			assert.True(t, success)
		})
		t.Run("Should support a 'deny-egress' rule policy for egress-node-peer", func(t *testing.T) {
			// harry-potter-1 is our client pod in gryffindor namespace
			// ensure egress is DENIED to rest of the nodes from gryffindor; egressRule at index2 should take effect
			success := kubernetes.PokeServer(t, s.ClientSet, &s.KubeConfig, "network-policy-conformance-gryffindor", "harry-potter-1", "tcp",
				serverPod.Status.PodIP, int32(36364), s.TimeoutConfig.RequestTimeout, false)
			assert.True(t, success)
			success = kubernetes.PokeServer(t, s.ClientSet, &s.KubeConfig, "network-policy-conformance-gryffindor", "harry-potter-1", "udp",
				serverPod.Status.PodIP, int32(34346), s.TimeoutConfig.RequestTimeout, false)
			assert.True(t, success)
			success = kubernetes.PokeServer(t, s.ClientSet, &s.KubeConfig, "network-policy-conformance-gryffindor", "harry-potter-1", "sctp",
				serverPod.Status.PodIP, int32(9003), s.TimeoutConfig.RequestTimeout, false)
			assert.True(t, success)
		})
	},
}

var AdminNetworkPolicyEgressInlineCIDRPeers = suite.ConformanceTest{
	ShortName:   "AdminNetworkPolicyEgressInlineCIDRPeers",
	Description: "Tests support for egress traffic to CIDR peers using admin network policy API based on a server and client model",
	Features: []suite.SupportedFeature{
		suite.SupportAdminNetworkPolicy,
		suite.SupportAdminNetworkPolicyEgressInlineCIDRPeers,
	},
	Manifests: []string{"base/admin_network_policy/experimental-egress-selector-rules.yaml"},
	Test: func(t *testing.T, s *suite.ConformanceTestSuite) {
		ctx, cancel := context.WithTimeout(context.Background(), s.TimeoutConfig.GetTimeout)
		defer cancel()
		// This test uses `node-and-cidr-as-peers-example` ANP
		t.Run("Should support a 'deny-egress' rule policy for egress-cidr-peer", func(t *testing.T) {
			// harry-potter-1 is our client pod in gryffindor namespace
			// Let us pick a pod in ravenclaw namespace and try to connect, it won't work
			// ensure egress is DENIED to 0.0.0.0/0 from gryffindor; egressRule at index2 should take effect
			// luna-lovegood-0 is our server pod in ravenclaw namespace
			serverPod := &v1.Pod{}
			err := s.Client.Get(ctx, client.ObjectKey{
				Namespace: "network-policy-conformance-ravenclaw",
				Name:      "luna-lovegood-0",
			}, serverPod)
			require.NoErrorf(t, err, "unable to fetch the server pod")
			success := kubernetes.PokeServer(t, s.ClientSet, &s.KubeConfig, "network-policy-conformance-gryffindor", "harry-potter-1", "tcp",
				serverPod.Status.PodIP, int32(80), s.TimeoutConfig.RequestTimeout, false)
			assert.True(t, success)
			success = kubernetes.PokeServer(t, s.ClientSet, &s.KubeConfig, "network-policy-conformance-gryffindor", "harry-potter-1", "udp",
				serverPod.Status.PodIP, int32(53), s.TimeoutConfig.RequestTimeout, false)
			assert.True(t, success)
			success = kubernetes.PokeServer(t, s.ClientSet, &s.KubeConfig, "network-policy-conformance-gryffindor", "harry-potter-1", "sctp",
				serverPod.Status.PodIP, int32(9003), s.TimeoutConfig.RequestTimeout, false)
			assert.True(t, success)
			// Let us pick a pod in hufflepuff namespace and try to connect, it won't work
			// ensure egress is DENIED to 0.0.0.0/0 from gryffindor; egressRule at index2 should take effect
			// cedric-diggory-0 is our server pod in hufflepuff namespace
			serverPod = &v1.Pod{}
			err = s.Client.Get(ctx, client.ObjectKey{
				Namespace: "network-policy-conformance-hufflepuff",
				Name:      "cedric-diggory-0",
			}, serverPod)
			require.NoErrorf(t, err, "unable to fetch the server pod")
			success = kubernetes.PokeServer(t, s.ClientSet, &s.KubeConfig, "network-policy-conformance-gryffindor", "harry-potter-1", "tcp",
				serverPod.Status.PodIP, int32(80), s.TimeoutConfig.RequestTimeout, false)
			assert.True(t, success)
			success = kubernetes.PokeServer(t, s.ClientSet, &s.KubeConfig, "network-policy-conformance-gryffindor", "harry-potter-1", "udp",
				serverPod.Status.PodIP, int32(53), s.TimeoutConfig.RequestTimeout, false)
			assert.True(t, success)
			success = kubernetes.PokeServer(t, s.ClientSet, &s.KubeConfig, "network-policy-conformance-gryffindor", "harry-potter-1", "sctp",
				serverPod.Status.PodIP, int32(9003), s.TimeoutConfig.RequestTimeout, false)
			assert.True(t, success)
		})
		// To test allow CIDR rule, insert the following rule at index0
		//- name: "allow-egress-to-specific-podIPs"
		//  action: "Allow"
		//  to:
		//  - networks:
		//	  - luna-lovegood-0.IP
		//    - cedric-diggory-0.IP
		t.Run("Should support an 'allow-egress' rule policy for egress-cidr-peer", func(t *testing.T) {
			serverPodRavenclaw := &v1.Pod{}
			err := s.Client.Get(ctx, client.ObjectKey{
				Namespace: "network-policy-conformance-ravenclaw",
				Name:      "luna-lovegood-0",
			}, serverPodRavenclaw)
			require.NoErrorf(t, err, "unable to fetch the server pod")
			serverPodHufflepuff := &v1.Pod{}
			err = s.Client.Get(ctx, client.ObjectKey{
				Namespace: "network-policy-conformance-hufflepuff",
				Name:      "cedric-diggory-0",
			}, serverPodHufflepuff)
			require.NoErrorf(t, err, "unable to fetch the server pod")
			anp := &v1alpha1.AdminNetworkPolicy{}
			err = s.Client.Get(ctx, client.ObjectKey{
				Name: "node-and-cidr-as-peers-example",
			}, anp)
			require.NoErrorf(t, err, "unable to fetch the admin network policy")
			mutate := anp.DeepCopy()
			var mask string
			if net.IsIPv4String(serverPodRavenclaw.Status.PodIP) {
				mask = "/32"
			} else {
				mask = "/128"
			}
			// insert new rule at index0; append the rest of the rules in the node-and-cidr-as-peers-example
			newRule := []v1alpha1.AdminNetworkPolicyEgressRule{
				{
					Name:   "allow-egress-to-specific-podIPs",
					Action: "Allow",
					To: []v1alpha1.AdminNetworkPolicyEgressPeer{
						{
							Networks: []v1alpha1.CIDR{
								v1alpha1.CIDR(serverPodRavenclaw.Status.PodIP + mask),
								v1alpha1.CIDR(serverPodHufflepuff.Status.PodIP + mask),
							},
						},
					},
				},
			}
			mutate.Spec.Egress = append(newRule, mutate.Spec.Egress...)
			err = s.Client.Patch(ctx, mutate, client.MergeFrom(anp))
			require.NoErrorf(t, err, "unable to patch the admin network policy")
			// harry-potter-0 is our client pod in gryffindor namespace
			// ensure egress is ALLOWED to luna-lovegood-0.IP and cedric-diggory-0.IP
			// new egressRule at index0 should take effect
			success := kubernetes.PokeServer(t, s.ClientSet, &s.KubeConfig, "network-policy-conformance-gryffindor", "harry-potter-1", "tcp",
				serverPodRavenclaw.Status.PodIP, int32(80), s.TimeoutConfig.RequestTimeout, true)
			assert.True(t, success)
			success = kubernetes.PokeServer(t, s.ClientSet, &s.KubeConfig, "network-policy-conformance-gryffindor", "harry-potter-1", "udp",
				serverPodRavenclaw.Status.PodIP, int32(53), s.TimeoutConfig.RequestTimeout, true)
			assert.True(t, success)
			success = kubernetes.PokeServer(t, s.ClientSet, &s.KubeConfig, "network-policy-conformance-gryffindor", "harry-potter-1", "sctp",
				serverPodRavenclaw.Status.PodIP, int32(9003), s.TimeoutConfig.RequestTimeout, true)
			assert.True(t, success)
			success = kubernetes.PokeServer(t, s.ClientSet, &s.KubeConfig, "network-policy-conformance-gryffindor", "harry-potter-1", "tcp",
				serverPodHufflepuff.Status.PodIP, int32(80), s.TimeoutConfig.RequestTimeout, true)
			assert.True(t, success)
			success = kubernetes.PokeServer(t, s.ClientSet, &s.KubeConfig, "network-policy-conformance-gryffindor", "harry-potter-1", "udp",
				serverPodHufflepuff.Status.PodIP, int32(53), s.TimeoutConfig.RequestTimeout, true)
			assert.True(t, success)
			success = kubernetes.PokeServer(t, s.ClientSet, &s.KubeConfig, "network-policy-conformance-gryffindor", "harry-potter-1", "sctp",
				serverPodHufflepuff.Status.PodIP, int32(9003), s.TimeoutConfig.RequestTimeout, true)
			assert.True(t, success)
		})
	},
}
