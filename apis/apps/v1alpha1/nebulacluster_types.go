/*
Copyright 2021 Vesoft Inc.

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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NebulaClusterConditionType represents a nebula cluster condition value.
type NebulaClusterConditionType string

const (
	// NebulaClusterReady indicates that the nebula cluster is ready or not.
	// This is defined as:
	// - All workloads are up to date (currentRevision == updateRevision).
	// - All nebula component pods are healthy.
	NebulaClusterReady NebulaClusterConditionType = "Ready"
)

// ComponentPhase is the current state of component
type ComponentPhase string

const (
	// RunningPhase represents normal state of nebula cluster.
	RunningPhase ComponentPhase = "Running"
	// UpgradePhase represents the upgrade state of nebula cluster.
	UpgradePhase ComponentPhase = "Upgrade"
	// ScaleInPhase represents the scaling state of nebula cluster.
	ScaleInPhase ComponentPhase = "ScaleIn"
	// ScaleOutPhase represents the scaling state of nebula cluster.
	ScaleOutPhase ComponentPhase = "ScaleOut"
	// UpdatePhase represents update state of nebula cluster.
	UpdatePhase ComponentPhase = "Update"
)

// NebulaClusterSpec defines the desired state of NebulaCluster
type NebulaClusterSpec struct {
	// graphd spec
	Graphd *GraphdSpec `json:"graphd"`

	// Metad spec
	Metad *MetadSpec `json:"metad"`

	// Storaged spec
	Storaged *StoragedSpec `json:"storaged"`

	// +optional
	Reference WorkloadReference `json:"reference,omitempty"`

	// +kubebuilder:default=default-scheduler
	// +optional
	SchedulerName string `json:"schedulerName"`

	// Flag to enable/disable pv reclaim while the nebula cluster deleted , default false
	// +optional
	EnablePVReclaim *bool `json:"enablePVReclaim,omitempty"`

	// +kubebuilder:default=IfNotPresent
	ImagePullPolicy *corev1.PullPolicy `json:"imagePullPolicy,omitempty"`

	// +optional
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`

	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// UpdatePolicy indicates how pods should be updated
	// +optional
	UpdatePolicy string `json:"strategy,omitempty"`
}

// NebulaClusterStatus defines the observed state of NebulaCluster
type NebulaClusterStatus struct {
	Graphd     ComponentStatus          `json:"graphd,omitempty"`
	Metad      ComponentStatus          `json:"metad,omitempty"`
	Storaged   ComponentStatus          `json:"storaged,omitempty"`
	Conditions []NebulaClusterCondition `json:"conditions,omitempty"`
}

// ComponentStatus is the status and version of a nebula component.
type ComponentStatus struct {
	Version  string         `json:"version,omitempty"`
	Phase    ComponentPhase `json:"phase,omitempty"`
	Workload WorkloadStatus `json:"workload,omitempty"`
}

// WorkloadStatus describes the status of a specified workload.
type WorkloadStatus struct {
	// ObservedGeneration is the most recent generation observed for this Workload. It corresponds to the
	// Workload's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// The number of ready replicas.
	// +optional
	ReadyReplicas int32 `json:"readyReplicas,omitempty"`

	// Replicas is the most recently observed number of replicas.
	Replicas int32 `json:"replicas"`

	// The number of pods in current version.
	UpdatedReplicas int32 `json:"updatedReplicas"`

	// The number of ready current revision replicas for this Workload.
	// +optional
	UpdatedReadyReplicas int32 `json:"updatedReadyReplicas,omitempty"`

	// Count of hash collisions for the Workload.
	// +optional
	CollisionCount *int32 `json:"collisionCount,omitempty"`

	// CurrentRevision, if not empty, indicates the current version of the Workload.
	CurrentRevision string `json:"currentRevision"`

	// updateRevision, if not empty, indicates the version of the Workload used to generate Pods in the sequence
	UpdateRevision string `json:"updateRevision,omitempty"`
}

// NebulaClusterCondition describes the state of a nebula cluster at a certain point.
type NebulaClusterCondition struct {
	// Type of the condition.
	Type NebulaClusterConditionType `json:"type"`
	// Status of the condition, one of True, False, Unknown.
	Status corev1.ConditionStatus `json:"status"`
	// The last time this condition was updated.
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`
	// Last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// The reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty"`
	// A human readable message indicating details about the transition.
	// +optional
	Message string `json:"message,omitempty"`
}

// A WorkloadReference refers to a CustomResourceDefinition by name.
type WorkloadReference struct {
	// Name of the referenced CustomResourceDefinition.
	// eg. statefulsets.apps
	Name string `json:"name"`

	// Version indicate which version should be used if CRD has multiple versions
	// by default it will use the first one if not specified
	Version string `json:"version,omitempty"`
}

// GraphdSpec defines the desired state of Graphd
type GraphdSpec struct {
	PodSpec `json:",inline"`

	// Config defines a graphd configuration load into ConfigMap
	Config map[string]string `json:"config,omitempty"`

	// Service defines a k8s service of Graphd cluster.
	// +optional
	Service *GraphdServiceSpec `json:"service,omitempty"`

	// K8S persistent volume claim for Graphd data storage.
	// +optional
	StorageClaim *StorageClaim `json:"storageClaim,omitempty"`
}

// MetadSpec defines the desired state of Metad
type MetadSpec struct {
	PodSpec `json:",inline"`

	// Config defines a metad configuration load into ConfigMap
	Config map[string]string `json:"config,omitempty"`

	// Service defines a Kubernetes service of Metad cluster.
	// +optional
	Service *ServiceSpec `json:"service,omitempty"`

	// K8S persistent volume claim for Metad data storage.
	// +optional
	StorageClaim *StorageClaim `json:"storageClaim,omitempty"`
}

// StoragedSpec defines the desired state of Storaged
type StoragedSpec struct {
	PodSpec `json:",inline"`

	// Config defines a storaged configuration load into ConfigMap
	Config map[string]string `json:"config,omitempty"`

	// Service defines a Kubernetes service of Storaged cluster.
	// +optional
	Service *ServiceSpec `json:"service,omitempty"`

	// K8S persistent volume claim for Storaged data storage.
	// +optional
	StorageClaim *StorageClaim `json:"storageClaim,omitempty"`
}

// PodSpec is a common set of k8s resource configs for nebula components.
type PodSpec struct {
	// K8S deployment replicas setting.
	// +kubebuilder:validation:Minimum=0
	Replicas *int32 `json:"replicas,omitempty"`

	// K8S resources settings.
	// +optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// Container environment variables.
	// +optional
	EnvVars []corev1.EnvVar `json:"env,omitempty"`

	// +kubebuilder:default=vesoft/graphd
	// +optional
	Image string `json:"image,omitempty"`

	// Version tag for docker images
	// +optional
	Version string `json:"version,omitempty"`

	// K8S pod annotations.
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// +optional
	Labels map[string]string `json:"labels,omitempty"`
}

// StorageClaim contains details of storages
type StorageClaim struct {
	// Resources represents the minimum resources the volume should have.
	// +optional
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`

	// Name of the StorageClass required by the claim.
	// +optional
	StorageClassName *string `json:"storageClassName,omitempty"`
}

// GraphdServiceSpec is the service spec of graphd
type GraphdServiceSpec struct {
	ServiceSpec `json:",inline"`

	// LoadBalancerIP is the loadBalancerIP of service
	// +optional
	LoadBalancerIP *string `json:"loadBalancerIP,omitempty"`

	// ExternalTrafficPolicy of the service
	// +optional
	ExternalTrafficPolicy *corev1.ServiceExternalTrafficPolicyType `json:"externalTrafficPolicy,omitempty"`
}

// ServiceSpec is a common set of k8s service configs.
type ServiceSpec struct {
	Type corev1.ServiceType `json:"type,omitempty"`

	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// +optional
	Selector map[string]string `json:"selector,omitempty"`

	// ClusterIP is the clusterIP of service
	// +optional
	ClusterIP *string `json:"clusterIP,omitempty"`

	// +optional
	PublishNotReadyAddresses bool `json:"publishNotReadyAddresses,omitempty"`
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=nc
// +kubebuilder:printcolumn:name="GRAPHD-DESIRED",type="string",JSONPath=".spec.graphd.replicas",description="The desired number of graphd pods."
// +kubebuilder:printcolumn:name="GRAPHD-READY",type="string",JSONPath=".status.graphd.workload.readyReplicas",description="The number of graphd pods ready."
// +kubebuilder:printcolumn:name="METAD-DESIRED",type="string",JSONPath=".spec.metad.replicas",description="The desired number of metad pods."
// +kubebuilder:printcolumn:name="METAD-READY",type="string",JSONPath=".status.metad.workload.readyReplicas",description="The number of metad pods ready."
// +kubebuilder:printcolumn:name="STORAGED-DESIRED",type="string",JSONPath=".spec.storaged.replicas",description="The desired number of storaged pods."
// +kubebuilder:printcolumn:name="STORAGED-READY",type="string",JSONPath=".status.storaged.workload.readyReplicas",description="The number of storaged pods ready."
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp",description="CreationTimestamp is a timestamp representing the server time when this object was created. It is represented in RFC3339 form and is in UTC."

// NebulaCluster is the Schema for the nebulaclusters API
type NebulaCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NebulaClusterSpec   `json:"spec,omitempty"`
	Status NebulaClusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// NebulaClusterList contains a list of NebulaCluster
type NebulaClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NebulaCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NebulaCluster{}, &NebulaClusterList{})
}
