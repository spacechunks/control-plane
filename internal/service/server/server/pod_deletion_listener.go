package server

import (
	"context"
	"log/slog"

	"github.com/spacechunks/control-plane/internal/service/server/kube"
	corev1 "k8s.io/api/core/v1"
)

type PodDeletionListener struct {
	logger     *slog.Logger
	serverRepo Repository
	clusterSet kube.ClusterSet
	namespace  string
	informers  []*kube.Informer
}

func NewPodDeletionListener(
	logger *slog.Logger,
	serverRepo Repository,
	clusterSet kube.ClusterSet,
	namespace string,
) *PodDeletionListener {
	return &PodDeletionListener{
		logger:     logger,
		serverRepo: serverRepo,
		clusterSet: clusterSet,
		namespace:  namespace,
	}
}

func (l PodDeletionListener) Start(ctx context.Context) error {
	informers := make([]*kube.Informer, 0, len(l.clusterSet))

	for clusterName, set := range l.clusterSet {
		i := kube.NewInformer(set, l.namespace, map[string]string{
			"chunks.cloud/server-type": "minigame",
		})

		i.AddPodDeleteHandler(func(pod *corev1.Pod) {
			l.podDeleteHandler(clusterName, pod)
		})

		informers = append(informers, i)

		if err := i.Start(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (l PodDeletionListener) Stop() {
	for _, i := range l.informers {
		i.Stop()
	}
}

func (l PodDeletionListener) podDeleteHandler(clusterName string, pod *corev1.Pod) {
	ctx := context.Background()

	id, ok := pod.Labels["chunks.cloud/server-id"]
	if !ok {
		l.logger.ErrorContext(
			ctx,
			"pod label chunks.cloud/server-id not found",
			"cluster_name", clusterName,
			"pod_name", pod.Name,
			"namespace", pod.Namespace,
		)
		return
	}

	l.logger.ErrorContext(ctx, "deleting server", "server_id", id)

	// FIXME: right now this thing has a big flaw: if this call errors, the pod is gone, but the db entry is still
	//        present. right now, we don't care about this, but this might need fixing at some point.
	if _, err := l.serverRepo.DeleteServer(ctx, id); err != nil {
		l.logger.ErrorContext(
			ctx,
			"failed to delete server",
			"err", err,
			"server_id", id,
		)
	}
}
