/*
Copyright 2023.

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
	"bytes"
	"context"
	"crypto/tls"
	"gopkg.in/gomail.v2"
	"html/template"
	coreV1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	pvmonitorV1 "pvmonitor/api/v1"
	"regexp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"time"
)

// PvMonitorReconciler reconciles a PvMonitor object
type PvMonitorReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	DynamicClient *dynamic.DynamicClient
}

// +kubebuilder:rbac:groups=xinyu.com,resources=pvmonitors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=xinyu.com,resources=pvmonitors/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=xinyu.com,resources=pvmonitors/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=pv,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// the PvMonitor object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *PvMonitorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("reconcile start", "namespace", req.Namespace, "resource name", req.Name)

	pvMonitor := pvmonitorV1.PvMonitor{}

	if err := r.Client.Get(ctx, req.NamespacedName, &pvMonitor); err != nil {
		logger.Error(err, "get resource err")
		return ctrl.Result{}, err
	}

	if pvMonitor.Status.Date == time.Now().Format("20060102") {
		logger.Info("has send")
		return ctrl.Result{}, nil
	}

	unStructList, err := r.DynamicClient.Resource(schema.GroupVersionResource{Group: "", Version: "v1", Resource: "persistentvolumes"}).
		List(ctx, v1.ListOptions{})
	if err != nil {
		logger.Error(err, "get pv list err")
		return ctrl.Result{}, err
	}

	pvs := make([]coreV1.PersistentVolume, 0)

	for _, v := range unStructList.Items {

		pv := new(coreV1.PersistentVolume)

		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(v.Object, pv); err != nil {
			logger.Error(err, "convert unstructed to persistent volume err")
			return ctrl.Result{}, err
		}
		if pvMonitor.Spec.Regex == "" {
			pvs = append(pvs, *pv)
			continue
		}
		matched, err := regexp.MatchString(pvMonitor.Spec.Regex, pv.Name)
		if err != nil {
			logger.Error(err, "regex match err")
			return ctrl.Result{}, err
		}
		if matched {
			pvs = append(pvs, *pv)
			continue
		}
	}
	tp, err := r.Template(pvs)
	if err != nil {
		logger.Error(err, "get template err")
		return ctrl.Result{}, err
	}

	if err := r.SendEmail(pvMonitor, tp); err != nil {
		logger.Error(err, "send email err")
		return ctrl.Result{}, err
	}

	pvMonitor.Status.Date = time.Now().Format("20060102")
	if err := r.Status().Update(ctx, &pvMonitor); err != nil {
		logger.Error(err, "update status error")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *PvMonitorReconciler) SendEmail(pvMonitor pvmonitorV1.PvMonitor, body string) error {

	m := gomail.NewMessage()
	m.SetAddressHeader("From", pvMonitor.Spec.Email.User, pvMonitor.Spec.Email.Subject)
	m.SetHeader("To", pvMonitor.Spec.Email.To...)
	if len(pvMonitor.Spec.Email.CC) > 0 {
		m.SetHeader("Cc", pvMonitor.Spec.Email.CC...)
	}
	m.SetHeader("Subject", pvMonitor.Spec.Email.Subject)

	m.SetBody("text/html", body)

	d := gomail.NewDialer(pvMonitor.Spec.Email.Host, pvMonitor.Spec.Email.Port, pvMonitor.Spec.Email.User, pvMonitor.Spec.Email.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func (r *PvMonitorReconciler) Template(pvs []coreV1.PersistentVolume) (string, error) {

	td := make([]TemplateData, 0)

	for _, v := range pvs {
		td = append(td, TemplateData{
			Name:    v.Name,
			Storage: v.Spec.Capacity.Storage().String(),
			Status:  string(v.Status.Phase),
		})
	}

	t, err := template.New("email").Parse(pvmonitorV1.Html)
	if err != nil {
		return "", err
	}
	writer := bytes.NewBuffer([]byte(""))

	if err := t.Execute(writer, struct {
		Monitors []TemplateData
	}{
		Monitors: td,
	}); err != nil {
		return "", err
	}
	return writer.String(), nil
}

type TemplateData struct {
	Name    string `json:"name"`
	Storage string `json:"storage"`
	Status  string `json:"status"`
}

// SetupWithManager sets up the controller with the Manager.
func (r *PvMonitorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&pvmonitorV1.PvMonitor{}).
		Complete(r)
}
