package framework

import(
	api "github.com/kubeware/messenger/apis/messenger/v1alpha1"
)

func (f *Invocation) NewMessagingService(
	name, namespace string, labels map[string]string,
	drive, secret string, to []string) *api.MessagingService {
	return &api.MessagingService{
		ObjectMeta: newObjectMeta(name, namespace, labels),
		Spec: api.MessagingServiceSpec{
			Drive: drive,
			To: to,
			CredentialSecretName: secret,
		},
	}
}