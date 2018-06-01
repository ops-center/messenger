package controller

import (
	"fmt"

	"github.com/appscode/kubernetes-webhook-util/admission"
	hooks "github.com/appscode/kubernetes-webhook-util/admission/v1beta1"
	webhook "github.com/appscode/kubernetes-webhook-util/admission/v1beta1/generic"
	"github.com/appscode/kutil/tools/queue"
	"github.com/golang/glog"
	"github.com/kubeware/messenger/apis/messenger"
	api "github.com/kubeware/messenger/apis/messenger/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
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
func (c *MessengerController) initNotifierWatcher() {
	c.notifierInformer = c.messengerInformerFactory.Messenger().V1alpha1().Notifiers().Informer()
	c.notifierQueue = queue.New(api.ResourceKindNotifier, c.MaxNumRequeues, c.NumThreads, c.reconcileNotifier)
	c.notifierInformer.AddEventHandler(queue.DefaultEventHandler(c.notifierQueue.GetQueue()))
	c.notifierLister = c.messengerInformerFactory.Messenger().V1alpha1().Notifiers().Lister()
}

func (c *MessengerController) reconcileNotifier(key string) error {
	obj, exist, err := c.notifierInformer.GetIndexer().GetByKey(key)
	if err != nil {
		glog.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}

	if !exist {
		glog.Warningf("Notifier %s does not exist anymore\n", key)
	} else {
		glog.Infof("Sync/Add/Update for Notifier %s\n", key)

		n := obj.(*api.Notifier)
		fmt.Println(n.Name)
	}
	return nil
}

func (c *MessengerController) deleteMessengerNotifier(repository *api.Notifier) error {
	return nil
}
