package postgres

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spacechunks/control-plane/internal/service/server/kube"
	"github.com/spacechunks/control-plane/internal/service/server/server"
)

func toDomainServer(r serverRelationsRow) (server.Server, error) {
	var labels map[string]string
	if err := json.Unmarshal(r.Server.Labels, &labels); err != nil {
		return server.Server{}, fmt.Errorf("unmarshal labels: %w", err)
	}

	ret := server.Server{
		ID:        r.Server.ID,
		Name:      r.Server.Name,
		Type:      server.ServerType(strings.ToLower(string(r.Server.Type))),
		Labels:    labels,
		State:     server.ServerState(strings.ToLower(string(r.Server.State))),
		CreatedAt: r.Server.CreatedAt.UTC(),
		UpdatedAt: r.Server.UpdatedAt.UTC(),
		Cluster: kube.Cluster{
			ID:        r.Cluster.ID,
			Name:      r.Cluster.Name,
			Slots:     r.Cluster.Slots,
			CreatedAt: r.Cluster.CreatedAt,
			UpdateAt:  r.Cluster.UpdatedAt,
		},
		ImageURL: r.Server.ImageUrl,
	}

	if r.Server.Ipv4 != nil {
		ret.Network = &server.Network{
			IP4:  *r.Server.Ipv4,
			Port: uint16(r.Server.Port.Int32),
		}
	}

	return ret, nil
}
