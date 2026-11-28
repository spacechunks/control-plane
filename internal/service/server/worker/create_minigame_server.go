package worker

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"time"

	"github.com/riverqueue/river"
	"github.com/spacechunks/control-plane/internal/eventually"
	"github.com/spacechunks/control-plane/internal/service/server/job"
	"github.com/spacechunks/control-plane/internal/service/server/kube"
	"github.com/spacechunks/control-plane/internal/service/server/server"
	corev1 "k8s.io/api/core/v1"
)

type CreateMinigameServerWorker struct {
	river.WorkerDefaults[job.CreateMinigameServer]
	logger         *slog.Logger
	serverRepo     server.Repository
	logTailTimeout time.Duration
	podPollTimeout time.Duration
	clusterSet     kube.ClusterSet
	namespace      string
}

func NewCreateMinigameServerWorker(
	logger *slog.Logger,
	serverRepo server.Repository,
	logTailTimeout time.Duration,
	podPollTimeout time.Duration,
	clusterSet kube.ClusterSet,
	namespace string,
) *CreateMinigameServerWorker {
	return &CreateMinigameServerWorker{
		logger:         logger,
		serverRepo:     serverRepo,
		logTailTimeout: logTailTimeout,
		podPollTimeout: podPollTimeout,
		clusterSet:     clusterSet,
		namespace:      namespace,
	}
}

func (w CreateMinigameServerWorker) Work(ctx context.Context, riverJob *river.Job[job.CreateMinigameServer]) error {
	serv, err := w.serverRepo.ServerByID(ctx, riverJob.Args.ServerID)
	if err != nil {
		return fmt.Errorf("server by id: %w", err)
	}

	w.logger.InfoContext(ctx, "creating minigame server", "server_id", serv.ID)

	var (
		labels = map[string]string{
			"chunks.cloud/server-id":   serv.ID,
			"chunks.cloud/server-type": serv.Type.String(),
		}
		jobArgs = riverJob.Args
	)

	clientset, ok := w.clusterSet[serv.Cluster.Name]
	if !ok {
		return fmt.Errorf("cluster %s not found in cluster set", serv.Cluster.Name)
	}

	kubeSvc := kube.NewClient(clientset, w.namespace)

	if err := kubeSvc.CreateJob(ctx, serv.ID, serv.ImageURL, labels); err != nil {
		return fmt.Errorf("create job: %w", err)
	}

	var pod corev1.Pod
	err = eventually.Run(ctx, w.podPollTimeout, 200*time.Millisecond, func(eCtx context.Context) error {
		pods, err := kubeSvc.PodByLabels(eCtx, map[string]string{
			"chunks.cloud/server-id": jobArgs.ServerID,
		})
		if err != nil {
			return fmt.Errorf("get pods: %w", err)
		}

		if len(pods.Items) == 0 {
			return fmt.Errorf("no pod found")
		}

		pod = pods.Items[0]
		return nil
	})
	if err != nil {
		return fmt.Errorf("eventually: %w", err)
	}

	if pod.Status.Phase == corev1.PodFailed {
		w.logger.ErrorContext(
			ctx,
			"canceling job, because pod is in failed state",
			"server_id", serv.ID,
			"pod_name", pod.Name,
		)
		if _, err := w.serverRepo.UpdateServer(ctx, server.UpdateServerDBArgs{
			ID:        serv.ID,
			State:     new(server.ServerStateError),
			UpdatedAt: time.Now(),
		}); err != nil {
			return fmt.Errorf("update state: %w", err)
		}

		return river.JobCancel(fmt.Errorf("pod in failed state"))
	}

	w.logger.InfoContext(
		ctx,
		"pod found for minigame server",
		"pod_name", pod.Name,
		"namespace", pod.Namespace,
		"server_id", serv.ID,
		"cluster_name", serv.Cluster.Name,
	)

	err = eventually.Run(ctx, w.logTailTimeout, 200*time.Millisecond, func(eCtx context.Context) error {
		stream, err := kubeSvc.TailLogs(eCtx, pod.Name, "mc")
		if err != nil {
			return fmt.Errorf("tail logs: %w", err)
		}

		defer stream.Close()
		paperServerReadyRegex := regexp.MustCompile(`Done \(\d+\.\d+s\)! For help, type "help"`)

		r := server.NewLogWatcher(stream, paperServerReadyRegex)

		if err := r.Run(eCtx); err != nil {
			return fmt.Errorf("log watcher: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("watch logs: %w", err)
	}

	w.logger.InfoContext(ctx, "minigame server ready", "server_id", serv.ID)

	if _, err := w.serverRepo.UpdateServer(ctx, server.UpdateServerDBArgs{
		ID:        serv.ID,
		State:     new(server.ServerStateRunning),
		UpdatedAt: time.Now(),
	}); err != nil {
		return fmt.Errorf("update state: %w", err)
	}

	return nil
}
