package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindNotifier = "Notifier"
	ResourceNotifier     = "notifier"
	ResourceNotifiers    = "notifiers"

	ResourceKindNotification = "Notification"
	ResourceNotification     = "notification"
	ResourceNotifications    = "notifications"
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

	// To whom notification will be sent
	To []string `json:"to,omitempty"`

	// How this notification will be sent
	Notifier string `json:"notifier,omitempty"`

	// Secret name to which credential data is provided to send notification
	CredentialSecretName string `json:"credentialSecretName,omitempty"`
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

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Notification struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              NotificationSpec   `json:"spec,omitempty"`
	Status            NotificationStatus `json:"status,omitempty"`
}

type NotificationSpec struct {
	Service string `json:"service,omitempty"`
	Message string `json:"message,omitempty"`
	Email   string `json:"email,omitempty"`
	Chat    string `json:"chat,omitempty"`
}

type NotificationStatus struct {
	SentTimestamp *metav1.Timestamp `json:"sentTimestamp,omitempty"`
	ErrorMessage  string            `json:"errorMessage,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type NotificationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of Notification TPR objects
	Items []Notification `json:"items,omitempty"`
}
