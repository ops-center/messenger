package v1alpha1

import (
	crdutils "github.com/appscode/kutil/apiextensions/v1beta1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	crd_api "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

var (
	EnableStatusSubresource bool
)

func (p MessagingService) CustomResourceDefinition() *crd_api.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Version:       SchemeGroupVersion.Version,
		Plural:        ResourceMessagingServices,
		Singular:      ResourceMessagingService,
		Kind:          ResourceKindMessagingService,
		ResourceScope: string(apiextensions.NamespaceScoped),
		Labels: crdutils.Labels{
			LabelsMap: map[string]string{"app": "messenger"},
		},
		SpecDefinitionName:      "github.com/appscode/messenger/apis/messenger/v1alpha1.MessagingService",
		EnableValidation:        true,
		GetOpenAPIDefinitions:   GetOpenAPIDefinitions,
		EnableStatusSubresource: EnableStatusSubresource,
	})
}

func (p MessagingService) IsValid() error {
	return nil
}

func (p Message) CustomResourceDefinition() *crd_api.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Version:       SchemeGroupVersion.Version,
		Plural:        ResourceMessages,
		Singular:      ResourceMessage,
		Kind:          ResourceKindMessage,
		ResourceScope: string(apiextensions.NamespaceScoped),
		Labels: crdutils.Labels{
			LabelsMap: map[string]string{"app": "messenger"},
		},
		SpecDefinitionName:      "github.com/appscode/messenger/apis/messenger/v1alpha1.Message",
		EnableValidation:        true,
		GetOpenAPIDefinitions:   GetOpenAPIDefinitions,
		EnableStatusSubresource: EnableStatusSubresource,
	})
}

func (p Message) IsValid() error {
	return nil
}
