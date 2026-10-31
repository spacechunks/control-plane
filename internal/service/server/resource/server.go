package resource

import "net/netip"

type ServerState string

const (
	ServerStatePending  = "pending"
	ServerStateCreating = "creating"
	ServerStateRunning  = "running"
	ServerStateDeleting = "deleting"
	ServerStateError    = "error"
)

type Server struct {
	ID      string
	Name    string
	Labels  map[string]string
	Network Network
	State   ServerState
}

type Network struct {
	IP4  *netip.Addr
	Port *uint16
}
