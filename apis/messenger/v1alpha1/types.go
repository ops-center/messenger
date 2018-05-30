package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindNotifier = "Notifier"
	ResourceNotifier     = "notifier"
	ResourceNotifiers    = "notifiers"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Notifier defines a Notifier database.
type Notifier struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              NotifierSpec   `json:"spec,omitempty"`
	Status            NotifierStatus `json:"status,omitempty"`
}

type NotifierSpec struct {
	// Number of instances to deploy for a Notifier database.
	Replicas *int32 `json:"replicas,omitempty"`
}

type NotifierStatus struct {
	CreationTime *metav1.Time `json:"creationTime,omitempty"`
	Reason       string       `json:"reason,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type NotifierList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of Notifier TPR objects
	Items []Notifier `json:"items,omitempty"`
}
