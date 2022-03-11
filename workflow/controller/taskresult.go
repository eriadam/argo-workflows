package controller

import (
	"time"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"

	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	wfextvv1alpha1 "github.com/argoproj/argo-workflows/v3/pkg/client/informers/externalversions/workflow/v1alpha1"
	"github.com/argoproj/argo-workflows/v3/workflow/common"
)

func (wfc *WorkflowController) newWorkflowTaskResultInformers() map[string]cache.SharedIndexInformer {
	informers := make(map[string]cache.SharedIndexInformer)
	for cluster, clientset := range wfc.wfclientsets {
		labelSelector := labels.NewSelector().
			Add(wfc.newProfileReq(cluster)).
			Add(wfc.newInstanceIDReq()).
			String()
		log.WithField("labelSelector", labelSelector).
			WithField("cluster", cluster).Info("Watching task results")
		informer := wfextvv1alpha1.NewFilteredWorkflowTaskResultInformer(
			clientset,
			wfc.GetManagedNamespace(),
			20*time.Minute,
			cache.Indexers{},
			func(options *metav1.ListOptions) {
				options.LabelSelector = labelSelector
			},
		)
		informer.AddEventHandler(
			cache.ResourceEventHandlerFuncs{
				AddFunc: func(new interface{}) {
					result := new.(*wfv1.WorkflowTaskResult)
					namespace := common.MetaWorkflowNamespace(result)
					workflow := result.Labels[common.LabelKeyWorkflow]
					wfc.wfQueue.AddRateLimited(namespace + "/" + workflow)
				},
				UpdateFunc: func(_, new interface{}) {
					result := new.(*wfv1.WorkflowTaskResult)
					namespace := common.MetaWorkflowNamespace(result)
					workflow := result.Labels[common.LabelKeyWorkflow]
					wfc.wfQueue.AddRateLimited(namespace + "/" + workflow)
				},
			})
		informers[cluster] = informer
	}
	return informers
}
