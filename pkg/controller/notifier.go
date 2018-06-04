package controller

import (
	"fmt"
	"strings"

	"github.com/appscode/envconfig"
	"github.com/appscode/go-notify"
	"github.com/appscode/go-notify/unified"
	"github.com/appscode/kubernetes-webhook-util/admission"
	hooks "github.com/appscode/kubernetes-webhook-util/admission/v1beta1"
	webhook "github.com/appscode/kubernetes-webhook-util/admission/v1beta1/generic"
	"github.com/appscode/kutil/tools/queue"
	"github.com/golang/glog"
	"github.com/kubeware/messenger/apis/messenger"
	api "github.com/kubeware/messenger/apis/messenger/v1alpha1"
	"github.com/tamalsaha/go-oneliners"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"github.com/appscode/kutil/meta"
	"time"
)

func (c *MessengerController) NewNotifierWebhook() hooks.AdmissionHook {
	return webhook.NewGenericWebhook(
		schema.GroupVersionResource{
			Group:    "admission.messenger.kubeware.io",
			Version:  "v1alpha1",
			Resource: api.ResourceNotifiers,
		},
		api.ResourceNotifier,
		[]string{messenger.GroupName},
		api.SchemeGroupVersion.WithKind(api.ResourceKindNotifier),
		nil,
		&admission.ResourceHandlerFuncs{
			CreateFunc: func(obj runtime.Object) (runtime.Object, error) {
				return nil, obj.(*api.Notifier).IsValid()
			},
			UpdateFunc: func(oldObj, newObj runtime.Object) (runtime.Object, error) {
				return nil, newObj.(*api.Notifier).IsValid()
			},
		},
	)
}
func (c *MessengerController) initNotificationWatcher() {
	c.notificationInformer = c.messengerInformerFactory.Messenger().V1alpha1().Notifications().Informer()
	c.notificationQueue = queue.New(api.ResourceKindNotification, c.MaxNumRequeues, c.NumThreads, c.reconcileNotification)
	c.notificationInformer.AddEventHandler(queue.DefaultEventHandler(c.notificationQueue.GetQueue()))
	c.notificationLister = c.messengerInformerFactory.Messenger().V1alpha1().Notifications().Lister()
}

func (c *MessengerController) reconcileNotification(key string) error {
	obj, exist, err := c.notificationInformer.GetIndexer().GetByKey(key)
	if err != nil {
		glog.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}

	if !exist {
		glog.Warningf("Notifier %s does not exist anymore\n", key)
	} else {
		glog.Infof("Sync/Add/Update for Notifier %s\n", key)

		n := obj.(*api.Notification)
		fmt.Println(">>>>>>>>>>>> notification crd obj name", n.Name)
		oneliners.PrettyJson(*n, "notificationCrdObj")
		err := c.send(n)
		if err != nil {
			n.Status.ErrorMessage = fmt.Sprintf("Sending error: %v", err)
			glog.Errorf(n.Status.ErrorMessage)
		} else {
			n.Status.SentTimestamp = &metav1.Timestamp{Seconds: int64(time.Now().Second())}
		}
	}
	return nil
}

func (c *MessengerController) deleteMessengerNotifier(repository *api.Notifier) error {
	return nil
}

func (c *MessengerController) send(notification *api.Notification) error {
	fmt.Println(">>>>>>>>>>>>> Send().......")
	notifierObj, err := c.messengerClient.MessengerV1alpha1().Notifiers(notification.Namespace).Get(notification.Spec.Service, metav1.GetOptions{})
	if err != nil {
		fmt.Println(">>>>>>>", )
		return err
	}

	oneliners.PrettyJson(*notifierObj, "notifierCrdObj")

	notifierCred, err := c.getLoader(notifierObj.Spec.CredentialSecretName)
	if err != nil {
		return err
	}

	notifier, err := unified.LoadVia(strings.ToLower(notifierObj.Spec.Notifier), notifierCred)
	if err != nil {
		return err
	}

	switch n := notifier.(type) {
	case notify.ByEmail:
		return n.To(notifierObj.Spec.To[0], notifierObj.Spec.To[1:]...).
			WithSubject(notification.Spec.Email).
			WithBody(notification.Spec.Message).
			WithNoTracking().
			Send()
	case notify.BySMS:
		return n.To(notifierObj.Spec.To[0], notifierObj.Spec.To[1:]...).
			WithBody(notification.Spec.Email).
			Send()
	case notify.ByChat:
		return n.To(notifierObj.Spec.To[0], notifierObj.Spec.To[1:]...).
			WithBody(notification.Spec.Chat).
			Send()
	case notify.ByPush:
		return n.To(notifierObj.Spec.To...).
			WithBody(notification.Spec.Chat).
			Send()
	}

	return nil
}

func (c *MessengerController) getLoader(credentialSecretName string) (envconfig.LoaderFunc, error) {
	if credentialSecretName == "" {
		return func(key string) (string, bool) {
			return "", false
		}, nil
	}
	cfg, err := c.kubeClient.CoreV1().
		Secrets(meta.Namespace()).
		Get(credentialSecretName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return func(key string) (value string, found bool) {
		var bytes []byte
		bytes, found = cfg.Data[key]
		value = string(bytes)
		return
	}, nil
}
