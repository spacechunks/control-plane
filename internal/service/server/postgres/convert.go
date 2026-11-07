package postgres

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spacechunks/control-plane/internal/service/server/postgres/query"
	"github.com/spacechunks/control-plane/internal/service/server/server"
)

func toDomainServer(s query.Server) (server.Server, error) {
	var labels map[string]string
	if err := json.Unmarshal(s.Labels, &labels); err != nil {
		return server.Server{}, fmt.Errorf("unmarshal labels: %w", err)
	}

	ret := server.Server{
		ID:        s.ID,
		Name:      s.Name,
		Type:      server.ServerType(strings.ToLower(string(s.Type))),
		Labels:    labels,
		State:     server.ServerState(strings.ToLower(string(s.State))),
		CreatedAt: s.CreatedAt.UTC(),
		UpdatedAt: s.UpdatedAt.UTC(),
	}

	if s.Ipv4 != nil {
		ret.Network = &server.Network{
			IP4:  *s.Ipv4,
			Port: uint16(s.Port.Int32),
		}
	}

	return ret, nil
}
