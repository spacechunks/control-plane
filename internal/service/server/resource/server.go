package resource

import (
	"net/netip"
	"time"
)

type ServerState string

const (
	ServerStatePending  = "pending"
	ServerStateCreating = "creating"
	ServerStateRunning  = "running"
	ServerStateDeleting = "deleting"
	ServerStateError    = "error"
)

type ServerType string

const (
	ServerTypeEphemeral  = "ephemeral"
	ServerTypePersistent = "persistent"
)

type Server struct {
	ID        string
	Name      string
	Type      ServerType
	Labels    map[string]string
	Network   Network
	State     ServerState
	CreatedAt time.Time
	UpdateAt  time.Time
}

type Network struct {
	IP4  *netip.Addr
	Port *uint16
}
