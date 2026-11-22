package server

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"time"
	"uuid"

	"github.com/spacechunks/control-plane/internal/eventually"
	"github.com/spacechunks/control-plane/internal/service/server/kube"
	corev1 "k8s.io/api/core/v1"
)

type CreateMinigameServerArgs struct {
	Name   string
	Labels map[string]string
	Image  string
	// flavor
	// image
}

type Service interface {
	CreateMinigameServer(ctx context.Context, args CreateMinigameServerArgs) (Server, error)
}

type Backend struct {
	logger      *slog.Logger
	serverRepo  Repository
	clusterRepo kube.ClusterRepository
	clusterSet  kube.ClusterSet
	namespace   string
}

func NewService(
	logger *slog.Logger,
	serverRepo Repository,
	clusterRepo kube.ClusterRepository,
	clusterSet kube.ClusterSet,
	namespace string,
) *Backend {
	return &Backend{
		logger:      logger,
		serverRepo:  serverRepo,
		clusterRepo: clusterRepo,
		clusterSet:  clusterSet,
		namespace:   namespace,
	}
}

func (s Backend) CreateMinigameServer(ctx context.Context, args CreateMinigameServerArgs) (Server, error) {
	cluster, err := s.clusterRepo.BestCluster(ctx)
	if err != nil {
		return Server{}, fmt.Errorf("best cluster: %w", err)
	}

	clientset, ok := s.clusterSet[cluster.Name]
	if !ok {
		return Server{}, fmt.Errorf("cluster %s not found in cluster set", cluster.Name)
	}

	kubeSvc := kube.NewClient(clientset, "default")

	dbArgs := CreateServerDBArgs{
		ID:        uuid.NewV7().String(),
		Name:      args.Name,
		ClusterID: cluster.ID,
		Labels:    args.Labels,
		Type:      ServerTypeMinigame,
		State:     ServerStateCreating,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	server, err := s.serverRepo.CreateServer(ctx, dbArgs)
	if err != nil {
		return Server{}, fmt.Errorf("create: %w", err)
	}

	// TODO: move everything below here into CreateServerJob

	labels := map[string]string{
		"chunks.cloud/server-id":   server.ID,
		"chunks.cloud/server-type": server.Type.String(),
	}

	if err := kubeSvc.CreateJob(ctx, server.ID, args.Image, labels); err != nil {
		return Server{}, fmt.Errorf("create job: %w", err)
	}

	var pod corev1.Pod
	err = eventually.Run(ctx, 10*time.Second, 200*time.Millisecond, func(eCtx context.Context) error {
		pods, err := kubeSvc.PodByLabels(eCtx, map[string]string{
			"chunks.cloud/server-id": server.ID,
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
		return Server{}, fmt.Errorf("eventually: %w", err)
	}

	// TODO: make configurable
	err = eventually.Run(ctx, 5*time.Minute, 200*time.Millisecond, func(eCtx context.Context) error {
		stream, err := kubeSvc.TailLogs(eCtx, pod.Name, "mc")
		if err != nil {
			return fmt.Errorf("tail logs: %w", err)
		}

		defer stream.Close()
		paperServerReadyRegex := regexp.MustCompile(`Done \(\d+\.\d+s\)! For help, type "help"`)

		r := NewLogWatcher(stream, paperServerReadyRegex)

		if err := r.Run(eCtx); err != nil {
			return fmt.Errorf("log watcher: %w", err)
		}

		return nil
	})

	if _, err := s.serverRepo.UpdateServer(ctx, UpdateServerDBArgs{
		ID:        server.ID,
		State:     new(ServerStateRunning),
		UpdatedAt: time.Now(),
	}); err != nil {
		return Server{}, fmt.Errorf("update state: %w", err)
	}

	return Server{}, nil
}
