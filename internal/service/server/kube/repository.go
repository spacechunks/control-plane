package kube

import (
	"context"
	"time"

	apierr "github.com/spacechunks/control-plane/internal/error"
)

var ErrNoSlotsAvailable = apierr.ResourceExhausted("no cluster available that can host the server")

type Cluster struct {
	ID        string
	Name      string
	Slots     int32
	CreatedAt time.Time
	UpdateAt  time.Time
}

type ClusterRepository interface {
	BestCluster(ctx context.Context) (Cluster, error)
}
