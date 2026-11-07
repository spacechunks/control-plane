package fixture

import (
	"net/netip"
	"time"

	"github.com/spacechunks/control-plane/internal/service/server/server"
)

func Server(mod ...func(s *server.Server)) server.Server {
	s := server.Server{
		ID:     "019532bb-2bd3-73ea-afc2-99f368a3eb97",
		Name:   "test-server",
		Type:   server.ServerTypeEphemeral,
		Labels: map[string]string{"k": "v"},
		Network: &server.Network{
			IP4:  netip.MustParseAddr("192.0.2.1"),
			Port: 25565,
		},
		State:     server.ServerStateRunning,
		CreatedAt: time.Date(2025, 2, 23, 13, 12, 15, 0, time.UTC),
		UpdatedAt: time.Date(2025, 2, 28, 10, 26, 0, 0, time.UTC),
	}

	for _, fn := range mod {
		fn(&s)
	}

	return s
}
