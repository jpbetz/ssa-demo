package main

import (
	"context"
	"os"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/jpbetz/ssademo/api/v1alpha1"
	"github.com/jpbetz/ssademo/api/v1alpha1/ac"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var scheme *runtime.Scheme

func init() {
	bld := ctrl.SchemeBuilder{}
	bld.GroupVersion.Group = "samplecontroller.k8s.io"
	bld.GroupVersion.Version = "v1alpha1"
	bld.Register(&v1alpha1.Foo{}, &v1alpha1.FooList{})

	var err error
	scheme, err = bld.Build()
	if err != nil {
		panic(err)
	}
}

func main() {
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))
	initLog := ctrl.Log.WithName("main")

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{Scheme: scheme})
	if err != nil {
		initLog.Error(err, "unable to create manager")
		os.Exit(1)
	}

	err = ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Foo{}).
		// Owns(blah), etc
		Complete(&rec{
			Client: mgr.GetClient(),
			log:    ctrl.Log.WithName("foo"),
		})

	if err != nil {
		initLog.Error(err, "unable to create controller")
		os.Exit(1)
	}

	initLog.Info("starting")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		initLog.Error(err, "unable to run manager")
		os.Exit(1)
	}
	initLog.Info("shutting down")
}

type rec struct {
	client.Client
	log logr.Logger
}

func (r *rec) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.log.WithValues("name", req.NamespacedName)
	log.Info("reconciling")

	var obj v1alpha1.Foo
	if err := r.Get(ctx, req.NamespacedName, &obj); err != nil {
		if client.IgnoreNotFound(err) == nil {
			return ctrl.Result{}, nil
		}
		log.Error(err, "unable to get requested object")
		return ctrl.Result{}, err
	}

	if obj.Spec.Replicas != nil && *obj.Spec.Replicas == 2 {
		return ctrl.Result{}, nil
	}

	fieldManager := client.FieldOwner("kubebuilder-controller")

	foo := ac.Foo(obj.Name, obj.Namespace).
		SetSpec(ac.FooSpec().
			SetReplicas(2),
		)

	if err := r.Patch(context.TODO(), foo, client.Apply, fieldManager, client.ForceOwnership); err != nil {
		panic(err)
	}

	return ctrl.Result{}, nil
}
