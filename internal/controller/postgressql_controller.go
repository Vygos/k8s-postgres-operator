/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"time"

	v1 "k8s.io/api/apps/v1"
	v12 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	databasev1alpha1 "github.com/vygos/postgres-operator/api/v1alpha1"
	"github.com/vygos/postgres-operator/internal/controller/pg"
	"github.com/vygos/postgres-operator/internal/glog"
)

// PostgresSQLReconciler reconciles a PostgresSQL object
type PostgresSQLReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=database.vygos.io,resources=postgressqls,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=database.vygos.io,resources=postgressqls/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=database.vygos.io,resources=postgressqls/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PostgresSQL object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *PostgresSQLReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	glog.Log(ctx).Info("Starting reconcile postgres controller")

	postgres := databasev1alpha1.PostgresSQL{}
	err := r.Get(ctx, req.NamespacedName, &postgres)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if postgres.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// to registering our finalizer.
		if !controllerutil.ContainsFinalizer(&postgres, databasev1alpha1.PostgresFinalizer) {
			controllerutil.AddFinalizer(&postgres, databasev1alpha1.PostgresFinalizer)
			if err := r.Update(ctx, &postgres); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		if controllerutil.ContainsFinalizer(&postgres, databasev1alpha1.PostgresFinalizer) {
			r.deleteResources(ctx, req, postgres)
		}

		controllerutil.RemoveFinalizer(&postgres, databasev1alpha1.PostgresFinalizer)
		if err := r.Update(ctx, &postgres); err != nil {
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil
	}

	return r.createResourceIfNotExist(ctx, req, postgres)
}

func (r *PostgresSQLReconciler) createResourceIfNotExist(
	ctx context.Context,
	req ctrl.Request,
	postgres databasev1alpha1.PostgresSQL,
) (ctrl.Result, error) {
	pgStatefulSet := v1.StatefulSet{}

	err := r.Get(ctx, req.NamespacedName, &pgStatefulSet)
	if err != nil {
		return pg.CreateIfNotFound(func() error {
			return r.createResource(ctx, &postgres, pg.CreateStatefulSet(postgres, req.Namespace))
		}, err)
	}

	svc := v12.Service{}
	err = r.Get(ctx, pg.BuildNamespacedSvcName(req.Namespace, postgres.Name), &svc)
	if err != nil {
		return pg.CreateIfNotFound(func() error {
			return r.createResource(ctx, &postgres, pg.CreateService(postgres, req.Namespace))
		}, err)
	}

	return ctrl.Result{RequeueAfter: time.Second}, nil
}

func (r *PostgresSQLReconciler) createResource(
	ctx context.Context,
	postgres *databasev1alpha1.PostgresSQL,
	obj client.Object,
) error {
	logger := log.FromContext(ctx)

	err := r.Create(ctx, obj)
	if err != nil {
		logger.Error(err, fmt.Sprintf("error while creating resource %s", obj.GetName()))
		postgres.Status.Failed()
		return r.Status().Update(ctx, postgres)
	}

	postgres.Status.Running()
	return r.Status().Update(ctx, postgres)
}

// SetupWithManager sets up the controller with the Manager.
func (r *PostgresSQLReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&databasev1alpha1.PostgresSQL{}).
		Owns(&v1.StatefulSet{}).
		Complete(r)
}

func (r *PostgresSQLReconciler) deleteResources(
	ctx context.Context,
	req ctrl.Request,
	postgres databasev1alpha1.PostgresSQL,
) {
	var (
		pgStatefulSet = &v1.StatefulSet{}
		svc           = &v12.Service{}
	)
	_ = r.Get(ctx, req.NamespacedName, pgStatefulSet)
	_ = r.Get(ctx, pg.BuildNamespacedSvcName(req.Namespace, postgres.Name), svc)

	if err := r.Delete(ctx, pgStatefulSet); err != nil {
		glog.Log(ctx).Error(err, "Error while trying to delete resource statefulset", pgStatefulSet)
	}

	if err := r.Delete(ctx, svc); err != nil {
		glog.Log(ctx).Error(err, "Error while trying to delete resource service", svc)
	}
}
