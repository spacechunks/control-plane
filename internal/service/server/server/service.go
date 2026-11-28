package server

import (
	"context"
	"fmt"
	"log/slog"
	"time"
	"uuid"

	"github.com/spacechunks/control-plane/internal/queue"
	"github.com/spacechunks/control-plane/internal/service/server/kube"
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
	q           *queue.RiverClient
}

func NewService(
	logger *slog.Logger,
	serverRepo Repository,
	clusterRepo kube.ClusterRepository,
	clusterSet kube.ClusterSet,
	namespace string,
	q *queue.RiverClient,
) *Backend {
	return &Backend{
		logger:      logger,
		serverRepo:  serverRepo,
		clusterRepo: clusterRepo,
		clusterSet:  clusterSet,
		namespace:   namespace,
		q:           q,
	}
}

func (s Backend) CreateMinigameServer(ctx context.Context, args CreateMinigameServerArgs) (Server, error) {
	cluster, err := s.clusterRepo.BestCluster(ctx)
	if err != nil {
		return Server{}, fmt.Errorf("best cluster: %w", err)
	}

	dbArgs := CreateServerDBArgs{
		ID:        uuid.NewV7().String(),
		Name:      args.Name,
		ClusterID: cluster.ID,
		Labels:    args.Labels,
		Type:      ServerTypeMinigame,
		State:     ServerStateCreating,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ImageURL:  args.Image,
	}

	server, err := s.serverRepo.CreateServer(ctx, dbArgs, s.q.InsertJobTX)
	if err != nil {
		return Server{}, fmt.Errorf("create: %w", err)
	}

	return server, nil
}
