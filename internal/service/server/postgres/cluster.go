package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/spacechunks/control-plane/internal/service/server/kube"
)

func (db DB) BestCluster(ctx context.Context) (kube.Cluster, error) {
	release, q, err := db.queries(ctx)
	if err != nil {
		return kube.Cluster{}, fmt.Errorf("queries: %w", err)
	}
	defer release()

	row, err := q.BestCluster(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return kube.Cluster{}, kube.ErrNoSlotsAvailable
		}
		return kube.Cluster{}, fmt.Errorf("best cluster: %w", err)
	}
	return kube.Cluster{
		ID:        row.ID,
		Name:      row.Name,
		Slots:     row.Slots,
		CreatedAt: row.CreatedAt,
		UpdateAt:  row.UpdatedAt,
	}, nil
}
