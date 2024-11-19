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

// AdminNetworkPolicyPortApplyConfiguration represents a declarative configuration of the AdminNetworkPolicyPort type for use
// with apply.
type AdminNetworkPolicyPortApplyConfiguration struct {
	PortNumber *PortApplyConfiguration      `json:"portNumber,omitempty"`
	NamedPort  *string                      `json:"namedPort,omitempty"`
	PortRange  *PortRangeApplyConfiguration `json:"portRange,omitempty"`
}

// AdminNetworkPolicyPortApplyConfiguration constructs a declarative configuration of the AdminNetworkPolicyPort type for use with
// apply.
func AdminNetworkPolicyPort() *AdminNetworkPolicyPortApplyConfiguration {
	return &AdminNetworkPolicyPortApplyConfiguration{}
}

// WithPortNumber sets the PortNumber field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the PortNumber field is set to the value of the last call.
func (b *AdminNetworkPolicyPortApplyConfiguration) WithPortNumber(value *PortApplyConfiguration) *AdminNetworkPolicyPortApplyConfiguration {
	b.PortNumber = value
	return b
}

// WithNamedPort sets the NamedPort field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the NamedPort field is set to the value of the last call.
func (b *AdminNetworkPolicyPortApplyConfiguration) WithNamedPort(value string) *AdminNetworkPolicyPortApplyConfiguration {
	b.NamedPort = &value
	return b
}

// WithPortRange sets the PortRange field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the PortRange field is set to the value of the last call.
func (b *AdminNetworkPolicyPortApplyConfiguration) WithPortRange(value *PortRangeApplyConfiguration) *AdminNetworkPolicyPortApplyConfiguration {
	b.PortRange = value
	return b
}
