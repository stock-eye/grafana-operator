/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either exgdess or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/linclaus/grafana-operator/pkg/grafana"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	grafanav1 "github.com/linclaus/grafana-operator/api/v1"
	"github.com/spf13/viper"
)

var LOG_FINALIZER = "grafanaDashboard"

// GrafanaDashboardReconciler reconciles a GrafanaDashboard object
type GrafanaDashboardReconciler struct {
	client.Client
	Log     logr.Logger
	Scheme  *runtime.Scheme
	Grafana grafana.Grafana
}

// +kubebuilder:rbac:groups=grafana.monitoring.io,resources=grafanadashboards,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=grafana.monitoring.io,resources=grafanadashboards/status,verbs=get;update;patch

func (r *GrafanaDashboardReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	_ = r.Log.WithValues("GrafanaDashboard", req.NamespacedName)

	gd := &grafanav1.GrafanaDashboard{}
	err := r.Get(ctx, req.NamespacedName, gd)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// GrafanaDashboard deleted
	if !gd.DeletionTimestamp.IsZero() {
		r.Log.V(1).Info("Deleting GrafanaDashboard")

		//delete dashboard
		err = r.Grafana.DeleteDashboard(gd.Spec.Title)
		if err != nil {
			r.updateStatus(gd, "Failed")
			return ctrl.Result{}, nil
		}

		//remove finalizer flag
		gd.Finalizers = removeString(gd.Finalizers, LOG_FINALIZER)
		if err = r.Update(ctx, gd); err != nil {
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
		r.Log.V(1).Info("GrafanaDashboard deleted")
		return ctrl.Result{}, nil
	}

	// Add finalizer
	if !containsString(gd.Finalizers, LOG_FINALIZER) {
		gd.Finalizers = append(gd.Finalizers, LOG_FINALIZER)
		if err = r.Update(ctx, gd); err != nil {
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
	}

	// GrafanaDashboard update
	if viper.GetString("FORCE_CREATE") == "true" || gd.Status.Status != "Successful" {
		r.Log.V(1).Info("Updating GrafanaDashboard")
		err = r.Grafana.UpsertDashboard(gd)
		if err != nil {
			r.updateStatus(gd, "Failed")
			return ctrl.Result{}, err
		}
		r.Log.V(1).Info("GrafanaDashboard updated")
		r.updateStatus(gd, "Successful")
	}

	return ctrl.Result{}, nil
}

func (r *GrafanaDashboardReconciler) updateStatus(gd *grafanav1.GrafanaDashboard, status string) {
	gd.Status.Status = status
	if status == "Failed" {
		rty := gd.Status.RetryTimes
		if rty < 100 {
			gd.Status.RetryTimes = rty + 1
		}
	}
	r.Update(context.Background(), gd)
}

func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func (r *GrafanaDashboardReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&grafanav1.GrafanaDashboard{}).
		Complete(r)
}
