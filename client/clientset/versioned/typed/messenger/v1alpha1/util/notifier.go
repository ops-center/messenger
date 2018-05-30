package util

import (
	"github.com/appscode/kutil"
	api "github.com/appscode/messenger/apis/messenger/v1alpha1"
	cs "github.com/appscode/messenger/client/clientset/versioned/typed/messenger/v1alpha1"
	"github.com/evanphx/json-patch"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
)

func CreateOrPatchNotifier(c cs.MessengerV1alpha1Interface, meta metav1.ObjectMeta, transform func(alert *api.Notifier) *api.Notifier) (*api.Notifier, kutil.VerbType, error) {
	cur, err := c.Notifiers(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		glog.V(3).Infof("Creating Notifier %s/%s.", meta.Namespace, meta.Name)
		out, err := c.Notifiers(meta.Namespace).Create(transform(&api.Notifier{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Notifier",
				APIVersion: api.SchemeGroupVersion.String(),
			},
			ObjectMeta: meta,
		}))
		return out, kutil.VerbCreated, err
	} else if err != nil {
		return nil, kutil.VerbUnchanged, err
	}
	return PatchNotifier(c, cur, transform)
}

func PatchNotifier(c cs.MessengerV1alpha1Interface, cur *api.Notifier, transform func(*api.Notifier) *api.Notifier) (*api.Notifier, kutil.VerbType, error) {
	return PatchNotifierObject(c, cur, transform(cur.DeepCopy()))
}

func PatchNotifierObject(c cs.MessengerV1alpha1Interface, cur, mod *api.Notifier) (*api.Notifier, kutil.VerbType, error) {
	curJson, err := json.Marshal(cur)
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}

	modJson, err := json.Marshal(mod)
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}

	patch, err := jsonpatch.CreateMergePatch(curJson, modJson)
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}
	if len(patch) == 0 || string(patch) == "{}" {
		return cur, kutil.VerbUnchanged, nil
	}
	glog.V(3).Infof("Patching Notifier %s/%s with %s.", cur.Namespace, cur.Name, string(patch))
	out, err := c.Notifiers(cur.Namespace).Patch(cur.Name, types.MergePatchType, patch)
	return out, kutil.VerbPatched, err
}

func TryUpdateNotifier(c cs.MessengerV1alpha1Interface, meta metav1.ObjectMeta, transform func(*api.Notifier) *api.Notifier) (result *api.Notifier, err error) {
	attempt := 0
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		cur, e2 := c.Notifiers(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {
			result, e2 = c.Notifiers(cur.Namespace).Update(transform(cur.DeepCopy()))
			return e2 == nil, nil
		}
		glog.Errorf("Attempt %d failed to update Notifier %s/%s due to %v.", attempt, cur.Namespace, cur.Name, e2)
		return false, nil
	})

	if err != nil {
		err = errors.Errorf("failed to update Notifier %s/%s after %d attempts due to %v", meta.Namespace, meta.Name, attempt, err)
	}
	return
}

func UpdateNotifierStatus(c cs.MessengerV1alpha1Interface, cur *api.Notifier, transform func(*api.NotifierStatus) *api.NotifierStatus, useSubresource ...bool) (*api.Notifier, error) {
	if len(useSubresource) > 1 {
		return nil, errors.Errorf("invalid value passed for useSubresource: %v", useSubresource)
	}

	mod := &api.Notifier{
		TypeMeta:   cur.TypeMeta,
		ObjectMeta: cur.ObjectMeta,
		Spec:       cur.Spec,
		Status:     *transform(cur.Status.DeepCopy()),
	}

	if len(useSubresource) == 1 && useSubresource[0] {
		return c.Notifiers(cur.Namespace).UpdateStatus(mod)
	}

	out, _, err := PatchNotifierObject(c, cur, mod)
	return out, err
}
