/*
Copyright 2023.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PvMonitorSpec defines the desired state of PvMonitor
type PvMonitorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of PvMonitor. Edit pvmonitor_types.go to remove/update
	Regex string         `json:"regex,omitempty"`
	Email PvMonitorEmail `json:"email,omitempty"`
}

type PvMonitorEmail struct {
	Host     string   `json:"host"`
	Port     int      `json:"port"`
	User     string   `json:"user"`
	Password string   `json:"password"`
	Subject  string   `json:"subject"`
	To       []string `json:"to"`
	CC       []string `json:"cc"`
}

// PvMonitorStatus defines the observed state of PvMonitor
type PvMonitorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Date string `json:"date,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=pvmonitors,scope=Cluster

// PvMonitor is the Schema for the pvmonitors API
type PvMonitor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PvMonitorSpec   `json:"spec,omitempty"`
	Status PvMonitorStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PvMonitorList contains a list of PvMonitor
type PvMonitorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PvMonitor `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PvMonitor{}, &PvMonitorList{})
}
