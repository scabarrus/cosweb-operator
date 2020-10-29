package cosweb

import (
	"context"

	examplev1 "cosweb-operator/pkg/apis/example/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"fmt"
)

var log = logf.Log.WithName("controller_cosweb")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Cosweb Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileCosweb{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("cosweb-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Cosweb
	err = c.Watch(&source.Kind{Type: &examplev1.Cosweb{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Cosweb
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &examplev1.Cosweb{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &examplev1.Cosweb{},
	})
	if err != nil {
		return err
	}
	return nil
}

// blank assignment to verify that ReconcileCosweb implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileCosweb{}

// ReconcileCosweb reconciles a Cosweb object
type ReconcileCosweb struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Cosweb object and makes changes based on the state read
// and what is in the Cosweb.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.


func (r *ReconcileCosweb) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Cosweb")
	// Fetch the Cosweb instance
	instance := &examplev1.Cosweb{}
	result,err:=r.ReconcilePod(request,instance)
	fmt.Println("=======================ICI=========================")
	result,err=r.ReconcileService(request,instance)
	if err != nil{
		return result,err
	}
	
	return reconcile.Result{}, nil
}
func (r *ReconcileCosweb) ReconcilePod(request reconcile.Request,instance *examplev1.Cosweb) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Pod")
	

	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Define a new Pod object
	pod := newPodForCR(instance)
	//fmt.Println("Pod : ",pod)
	// Set Cosweb instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Pod already exists
	found := &corev1.Pod{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
		err = r.client.Create(context.TODO(), pod)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Pod created successfully - don't requeue
		//return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}
	// Pod already exists - don't requeue
	reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)
	return reconcile.Result{}, err
}
func (r *ReconcileCosweb) ReconcileService(request reconcile.Request,instance *examplev1.Cosweb) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Service Cosweb")
	
	// Check if this Pod already exists
	found := &corev1.Service{}
	//look if service with name instance.Name+"-svc" exist
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name+"-svc", Namespace: instance.Namespace}, found)
	//if service not exist for pod 
	if err != nil && errors.IsNotFound(err) {
		// Define a new Service object
		service := newServiceForPod(instance)
		fmt.Println("SERVICE : ",service)
		// Set Cosweb instance as the owner and controller
		if err := controllerutil.SetControllerReference(instance, service, r.scheme); err != nil {
			return reconcile.Result{}, err
		}
		reqLogger.Info("Creating service for Pod", "Pod.Namespace",instance.Namespace, "Pod.Name", instance.Name)
		err = r.client.Create(context.TODO(), service)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Service created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}
	// Pod already exists - don't requeue
	reqLogger.Info("Skip reconcile: Service already exists", "Service.Namespace", found.Namespace, "Service.Name", found.Name)
	
	return reconcile.Result{}, nil
}

func newServiceForPod(cr *examplev1.Cosweb) *corev1.Service{
	labels := cr.Spec.Labels
	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name+"-svc",
			Namespace: cr.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Protocol:   corev1.ProtocolTCP,
				Port:       8000,
				TargetPort: intstr.FromInt(8000),
				NodePort:   30685,
			}},
			Type: corev1.ServiceTypeNodePort,
		},
	}

	
	return s
}
// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(cr *examplev1.Cosweb) *corev1.Pod {
	labels := cr.Spec.Labels
	//- name: postgresCon
	//value: user=admin password=admin dbname=golang host=pg-golang-699c7cdc88-k8j4j.default.svc.cluster.local
	//port=5432 sslmode=disable

	var envVars  []corev1.EnvVar
	env :=corev1.EnvVar{"postgrescon","user=admin password=admin dbname=golang host=pg-golang-699c7cdc88-k8j4j.default.svc.cluster.local port=5432 sslmode=disable",nil}
	envVars=append(envVars,env)
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    cr.Spec.Name,
					Image:   cr.Spec.Image,
					Env: envVars,
					//Command: []string{"sleep", "3600"},
				},
			},
		},
	}
}
