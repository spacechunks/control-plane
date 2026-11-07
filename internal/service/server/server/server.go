package server

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
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Network struct {
	IP4  netip.Addr
	Port uint16
}
