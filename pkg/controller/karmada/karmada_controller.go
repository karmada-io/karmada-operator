package karmada

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	operatorapi "github.com/karmada-io/karmada-operator/pkg/apis/operator/v1alpha1"
)

const (
	// ControllerName is the controller name that will be used when reporting events.
	ControllerName = "karmada-controller"

	// name of the karmada controller finalizer
	KarmadaControllerFinalizerName = "operator.karmada.io/finalizer"
)

type KarmadaController struct {
	client.Client

	EventRecorder record.EventRecorder
}

// Reconcile performs a full reconciliation for the object referred to by the Request.
// The Controller will requeue the Request to be processed again if an error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (ctrl *KarmadaController) Reconcile(ctx context.Context, req controllerruntime.Request) (controllerruntime.Result, error) {
	startTime := time.Now()
	klog.V(4).InfoS("Started syncing karmada", "karmada", req, "startTime", startTime)
	defer func() {
		klog.V(4).InfoS("Finished syncing karmada", "karmada", req, "duration", time.Since(startTime))
	}()

	karmada := &operatorapi.Karmada{}
	if err := ctrl.Get(ctx, req.NamespacedName, karmada); err != nil {
		// The resource may no longer exist, in which case we stop processing.
		if errors.IsNotFound(err) {
			klog.V(2).InfoS("Karmada has been deleted", "karmada", req)
			return controllerruntime.Result{}, nil
		}
		return controllerruntime.Result{}, err
	}

	// Deep-copy otherwise we are mutating our cache.
	// TODO: Deep-copy only when needed.
	karmada = karmada.DeepCopy()

	// examine DeletionTimestamp to determine if object is under deletion
	if karmada.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !controllerutil.ContainsFinalizer(karmada, KarmadaControllerFinalizerName) {
			controllerutil.AddFinalizer(karmada, KarmadaControllerFinalizerName)
			if err := ctrl.Update(ctx, karmada); err != nil {
				return controllerruntime.Result{}, err
			}
		}
	} else {
		// The object is being deleted
		if controllerutil.ContainsFinalizer(karmada, KarmadaControllerFinalizerName) {
			// our finalizer is present, so lets handle any external dependency
			if err := ctrl.deleteUnableGCResources(karmada); err != nil {
				// if fail to delete the external dependency here, return with error
				// so that it can be retried
				return controllerruntime.Result{}, err
			}

			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(karmada, KarmadaControllerFinalizerName)
			if err := ctrl.Update(ctx, karmada); err != nil {
				return controllerruntime.Result{}, err
			}
			// Stop reconciliation as the item is being deleted
			return controllerruntime.Result{}, nil
		}
	}

	klog.V(2).InfoS("Reconciling karmada", "name", req.Name)

	// do reconcile

	return controllerruntime.Result{}, nil
}

func (ctrl *KarmadaController) deleteUnableGCResources(karmada *operatorapi.Karmada) error {
	klog.InfoS("Deleting unable gc resources", "karmada", klog.KObj(karmada))
	return nil
}

// SetupWithManager creates a controller and register to controller manager.
func (c *KarmadaController) SetupWithManager(mgr controllerruntime.Manager) error {
	return controllerruntime.NewControllerManagedBy(mgr).For(&operatorapi.Karmada{}).Complete(c)
}
