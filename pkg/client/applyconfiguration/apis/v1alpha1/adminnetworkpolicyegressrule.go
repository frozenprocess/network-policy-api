/*
Copyright The Kubernetes Authors.

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

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "sigs.k8s.io/network-policy-api/apis/v1alpha1"
)

// AdminNetworkPolicyEgressRuleApplyConfiguration represents a declarative configuration of the AdminNetworkPolicyEgressRule type for use
// with apply.
type AdminNetworkPolicyEgressRuleApplyConfiguration struct {
	Name   *string                                          `json:"name,omitempty"`
	Action *v1alpha1.AdminNetworkPolicyRuleAction           `json:"action,omitempty"`
	To     []AdminNetworkPolicyEgressPeerApplyConfiguration `json:"to,omitempty"`
	Ports  *[]AdminNetworkPolicyPortApplyConfiguration      `json:"ports,omitempty"`
}

// AdminNetworkPolicyEgressRuleApplyConfiguration constructs a declarative configuration of the AdminNetworkPolicyEgressRule type for use with
// apply.
func AdminNetworkPolicyEgressRule() *AdminNetworkPolicyEgressRuleApplyConfiguration {
	return &AdminNetworkPolicyEgressRuleApplyConfiguration{}
}

// WithName sets the Name field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Name field is set to the value of the last call.
func (b *AdminNetworkPolicyEgressRuleApplyConfiguration) WithName(value string) *AdminNetworkPolicyEgressRuleApplyConfiguration {
	b.Name = &value
	return b
}

// WithAction sets the Action field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Action field is set to the value of the last call.
func (b *AdminNetworkPolicyEgressRuleApplyConfiguration) WithAction(value v1alpha1.AdminNetworkPolicyRuleAction) *AdminNetworkPolicyEgressRuleApplyConfiguration {
	b.Action = &value
	return b
}

// WithTo adds the given value to the To field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the To field.
func (b *AdminNetworkPolicyEgressRuleApplyConfiguration) WithTo(values ...*AdminNetworkPolicyEgressPeerApplyConfiguration) *AdminNetworkPolicyEgressRuleApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithTo")
		}
		b.To = append(b.To, *values[i])
	}
	return b
}

func (b *AdminNetworkPolicyEgressRuleApplyConfiguration) ensureAdminNetworkPolicyPortApplyConfigurationExists() {
	if b.Ports == nil {
		b.Ports = &[]AdminNetworkPolicyPortApplyConfiguration{}
	}
}

// WithPorts adds the given value to the Ports field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Ports field.
func (b *AdminNetworkPolicyEgressRuleApplyConfiguration) WithPorts(values ...*AdminNetworkPolicyPortApplyConfiguration) *AdminNetworkPolicyEgressRuleApplyConfiguration {
	b.ensureAdminNetworkPolicyPortApplyConfigurationExists()
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithPorts")
		}
		*b.Ports = append(*b.Ports, *values[i])
	}
	return b
}
