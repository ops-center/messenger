package v1alpha1

import (
	crdutils "github.com/appscode/kutil/apiextensions/v1beta1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	crd_api "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

func (p Notifier) CustomResourceDefinition() *crd_api.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Version:       SchemeGroupVersion.Version,
		Plural:        ResourceNotifiers,
		Singular:      ResourceNotifier,
		Kind:          ResourceKindNotifier,
		ResourceScope: string(apiextensions.NamespaceScoped),
		Labels: crdutils.Labels{
			LabelsMap: map[string]string{"app": "messenger"},
		},
		SpecDefinitionName:    "github.com/kubeware/messenger/apis/messenger/v1alpha1.Notifier",
		EnableValidation:      true,
		GetOpenAPIDefinitions: GetOpenAPIDefinitions,
	})
}

func (p Notifier) IsValid() error {
	return nil
}
