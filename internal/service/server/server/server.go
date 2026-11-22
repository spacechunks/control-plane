package server

import (
	"net/netip"
	"time"

	"github.com/spacechunks/control-plane/internal/service/server/kube"
)

type ServerState string

const (
	ServerStatePending  ServerState = "pending"
	ServerStateCreating ServerState = "creating"
	ServerStateRunning  ServerState = "running"
	ServerStateDeleting ServerState = "deleting"
	ServerStateError    ServerState = "error"
)

func (s ServerState) String() string {
	return string(s)
}

type ServerType string

const (
	ServerTypeEphemeral  = "ephemeral"
	ServerTypePersistent = "persistent"
	ServerTypeMinigame   = "minigame"
)

func (t ServerType) String() string {
	return string(t)
}

type Server struct {
	ID        string
	Name      string
	Type      ServerType
	Labels    map[string]string
	Network   *Network
	State     ServerState
	Cluster   kube.Cluster
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Network struct {
	IP4  netip.Addr
	Port uint16
}
