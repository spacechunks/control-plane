package kube

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

type PodDeleteHandler func(pod *corev1.Pod)

type Informer struct {
	logger         *slog.Logger
	clientset      *kubernetes.Clientset
	namespace      string
	labelSelector  map[string]string
	deleteHandlers []PodDeleteHandler
	mu             sync.Mutex
	factory        informers.SharedInformerFactory
}

func NewInformer(clientset *kubernetes.Clientset, namespace string, labelSelector map[string]string) *Informer {
	return &Informer{
		clientset:     clientset,
		namespace:     namespace,
		labelSelector: labelSelector,
	}
}

func (i *Informer) AddPodDeleteHandler(handlers ...PodDeleteHandler) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if i.factory != nil {
		return
	}

	i.deleteHandlers = handlers
}

func (i *Informer) Start(ctx context.Context) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	factory := informers.NewSharedInformerFactoryWithOptions(
		i.clientset,
		0,
		informers.WithNamespace(i.namespace),
		informers.WithTweakListOptions(func(o *metav1.ListOptions) {
			o.LabelSelector = labels.SelectorFromSet(i.labelSelector).String()
		}),
	)

	podInformer := factory.Core().V1().Pods().TypedInformer()

	_, err := podInformer.AddTypedEventHandler(cache.TypedResourceEventHandlerFuncs[*corev1.Pod]{
		DeleteFunc: i.handlePodDelete,
	})
	if err != nil {
		return fmt.Errorf("add event handler: %w", err)
	}

	factory.Start(ctx.Done())

	if !cache.WaitForCacheSync(ctx.Done(), podInformer.HasSynced) {
		return fmt.Errorf("failed to wait for cache sync")
	}

	i.factory = factory
	return nil
}

func (i *Informer) Stop() {
	i.factory.Shutdown()
}

func (i *Informer) handlePodDelete(deleted cache.DeletedObject[*corev1.Pod]) {
	pod := deleted.OptionalObj
	// if nil, we missed the delete and the final state could
	// not be recovered; namespace/name are always available.
	if pod == nil {
		i.logger.Warn(
			"pod deleted, but no pod found",
			"pod_name", deleted.GetName(),
			"namespace", deleted.GetNamespace(),
		)
		return
	}

	for _, h := range i.deleteHandlers {
		h(pod)
	}
}
