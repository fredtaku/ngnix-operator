/*
Copyright 2025.

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

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	operatorv1alpha1 "github.com/fredtaku/ngnix-operator/api/v1alpha1"
)

func nginxDeployment() *appsv1.Deployment {
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginxDeployment",
			Namespace: "nginxDeploymentNS",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: pointer.Int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "nginx",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "nginx",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:    "nginx",
						Image:   "nginx:latest",
						Command: []string{"nginx"},
						Ports: []corev1.ContainerPort{{
							ContainerPort: 8080,
							Name:          "nginx",
						}},
					}},
				},
			},
		},
	}
	return dep
}

// MyNginxOperatorReconciler reconciles a MyNginxOperator object
type MyNginxOperatorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=operator.example.com,resources=mynginxoperators,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=operator.example.com,resources=mynginxoperators/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=operator.example.com,resources=mynginxoperators/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MyNginxOperator object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile

func (r *MyNginxOperatorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Step 1: Create the Deployment object
	deployment := nginxDeployment()

	// Step 2: Attempt to create it in the cluster
	err := r.Client.Create(ctx, deployment)
	if err != nil {
		if apierrors.IsAlreadyExists(err) {
			log.Info("Deployment already exists, skipping creation")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to create Deployment")
		return ctrl.Result{}, err
	}

	log.Info("Successfully created Deployment")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MyNginxOperatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&operatorv1alpha1.MyNginxOperator{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
