package main

import (
	"context"
	"time"

	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	runinformer "github.com/tektoncd/pipeline/pkg/client/injection/informers/pipeline/v1alpha1/run"
	runreconciler "github.com/tektoncd/pipeline/pkg/client/injection/reconciler/pipeline/v1alpha1/run"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/injection/sharedmain"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/reconciler"
)

const (
	controllerName = "replaceme-task-controller"
)

func main() {
	sharedmain.Main(controllerName, newController)
}

func newController(ctx context.Context, cmw configmap.Watcher) *controller.Impl {
	c := &Reconciler{
		// TODO: Provide fields that are useful when reconciling the
		// state of your Run (e.g., API clients).
	}
	impl := runreconciler.NewImpl(ctx, c, func(impl *controller.Impl) controller.Options {
		return controller.Options{
			AgentName: controllerName,
		}
	})
	r.enqueueAfter = impl.EnqueueAfter

	runinformer.Get(ctx).Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: FilterRunRef("my.api-group.dev/v0alpha0", "ReplaceMe"),
		Handler:    controller.HandleAll(impl.Enqueue),
	})

	return impl
}

// FilterRunRef returns a filter that can be passed to a Run Informer, which
// filters out Runs for apiVersion and kinds that this controller doesn't care
// about.
// TODO: Provide this as a helper function.
func FilterRunRef(apiVersion, kind string) func(interface{}) bool {
	return func(obj interface{}) bool {
		r, ok := obj.(*v1alpha1.Run)
		if !ok {
			// Somehow got informed of a non-Run object.
			// Ignore.
			return false
		}
		if r == nil || r.Spec.Ref == nil {
			// These are invalid, but just in case they get
			// created somehow, don't panic.
			return false
		}

		return r.Spec.Ref.APIVersion == apiVersion && r.Spec.Ref.Kind == v1alpha1.TaskKind(kind)
	}
}

type Reconciler struct {
	// TODO: Add fields that are useful to reconciling the state of your
	// Run, and provide them above.

	// enqueueAfter can be called from your reconciler to schedule a
	// reconciliation at a time in the future. This is not guaranteed to
	// happen at exactly the specified time.
	enqueueAfter func(*v1alpha1.Run, time.Duration)
}

// ReconcileKind implements Interface.ReconcileKind.
func (c *Reconciler) ReconcileKind(ctx context.Context, r *v1alpha1.Run) reconciler.Event {
	logger := logging.FromContext(ctx)
	logger.Infof("Reconciling %s/%s", r.Namespace, r.Name)

	// NB: Make sure to ignore reconciliations on done Runs; this can
	// happen when a global reconcile is triggered, and usually shouldn't
	// trigger any side effects.
	if r.IsDone() {
		return nil
	}

	// TODO: Update the Run to report that it's ongoing, or finished
	// successfully or unsuccessfully.

	return reconciler.NewEvent(corev1.EventTypeNormal, "RunReconciled", "Run reconciled: \"%s/%s\"", r.Namespace, r.Name)
}
