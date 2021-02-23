//go:generate go run sigs.k8s.io/controller-tools/cmd/controller-gen paths=. object crd
//go:generate go run sigs.k8s.io/controller-tools/cmd/controller-gen apply paths="./..."

// +groupName=samplecontroller.k8s.io
// +versionName=v1alpha1
package main

import (
	"context"
	"os"

	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	// corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/jpbetz/ssademo/ac"
)

var scheme *runtime.Scheme

func init() {
	bld := ctrl.SchemeBuilder{}
	bld.GroupVersion.Group = "samplecontroller.k8s.io"
	bld.GroupVersion.Version = "v1alpha1"
	bld.Register(&Foo{}, &FooList{})

	var err error
	scheme, err = bld.Build()
	if err != nil {
		panic(err)
	}

	// use `corev1.AddToScheme(scheme)` if you want to interact with core/v1,
	// for example
}

// +kubebuilder:object:generate=true
// +kubebuilder:object:root=true

type Foo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   FooSpec   `json:"spec"`
	Status FooStatus `json:"status"`
}

type (
	FooSpec struct {
		DeploymentName string `json:"deploymentName"`
		Replicas       *int32 `json:"replicas"`
	}
	FooStatus struct {
		AvailableReplicas int32 `json:"availableReplicas"`
	}
)

// +kubebuilder:object:root=true

type FooList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Foo `json:"items"`
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
		For(&Foo{}).
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

	var obj Foo
	if err := r.Get(ctx, req.NamespacedName, &obj); err != nil {
		if client.IgnoreNotFound(err) == nil {
			return ctrl.Result{}, nil
		}
		log.Error(err, "unable to get requested object")
		return ctrl.Result{}, err
	}

	var fieldManager client.FieldOwner = "kubebuilder-controller"
	foo := ac.Foo().
		SetSpec(ac.FooSpec().
			SetReplicas(2),
		)
	if err := r.Patch(context.TODO(), foo, client.Apply, fieldManager); err != nil {
		panic(err)
	}

	return ctrl.Result{}, nil
}
