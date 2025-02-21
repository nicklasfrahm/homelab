/*
Copyright 2025.

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

package v1beta1

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE! THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required. Any new fields you add
// must have json tags for the fields to be serialized.

// MAC is a MAC address that can be marshaled
// to and unmarshaled from YAML.
type MAC net.HardwareAddr

// String returns the string representation of a MAC address.
func (m MAC) String() string {
	builder := strings.Builder{}

	for index, octet := range m {
		if index > 0 {
			builder.WriteRune(':')
		}

		builder.WriteString(hex.EncodeToString([]byte{octet}))
	}

	return builder.String()
}

// UnmarshalYAML unmarshals a MAC address from a string.
func (m *MAC) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var mac string

	if err := unmarshal(&mac); err != nil {
		return fmt.Errorf("failed to unmarshal MAC address: %w", err)
	}

	hw, err := net.ParseMAC(mac)
	if err != nil {
		return fmt.Errorf("failed to parse MAC address: %w", err)
	}

	*m = MAC(hw)

	return nil
}

// MarshalYAML marshals a MAC address to a string.
func (m MAC) MarshalYAML() (interface{}, error) {
	return m.String(), nil
}

// UnmarshalJSON unmarshals a MAC address from a string.
func (m *MAC) UnmarshalJSON(data []byte) error {
	var mac string

	if err := json.Unmarshal(data, &mac); err != nil {
		return fmt.Errorf("failed to unmarshal MAC address: %w", err)
	}

	hw, err := net.ParseMAC(mac)
	if err != nil {
		return fmt.Errorf("failed to parse MAC address: %w", err)
	}

	*m = MAC(hw)

	return nil
}

// MarshalJSON marshals a MAC address to a string.
func (m MAC) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, m.String())), nil
}

// Interface describes a network interface of a Machine.
type Interface struct {
	// MAC is the MAC address of the interface.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Format=mac
	MAC MAC `json:"mac"`
}

// MachineSpecHardware defines the hardware configuration of a Machine.
type MachineSpecHardware struct {
	// Vendor is the manufacturer of the machine.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	Vendor string `json:"vendor"`
	// Model is the model of the machine.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	Model string `json:"model"`
}

// MachineSpec defines the desired state of a Machine.
type MachineSpec struct {
	// Hardware is the hardware configuration of the machine.
	// +kubebuilder:validation:Required
	Hardware MachineSpecHardware `json:"hardware"`
	// Interfaces describes the network interfaces of the machine.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	Interfaces []Interface `json:"interfaces"`
}

// MachineStatus defines the observed state of a Machine.
type MachineStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Machine defines a physical asset that can be used to provision infrastructure.
type Machine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MachineSpec   `json:"spec,omitempty"`
	Status MachineStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MachineList contains a list of Machine
type MachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Machine `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Machine{}, &MachineList{})
}
