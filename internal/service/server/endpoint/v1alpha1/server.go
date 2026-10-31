package v1alpha1

import (
	"context"
	"fmt"

	serverv1alpha1 "github.com/spacechunks/control-plane/api/server/v1alpha1"
	"github.com/spacechunks/control-plane/internal/ptr"
	"github.com/spacechunks/control-plane/internal/service/server"
)

func (s Handler) ListServers(ctx context.Context, req serverv1alpha1.ListServersRequestObject) (serverv1alpha1.ListServersResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s Handler) CreateServer(
	ctx context.Context,
	req serverv1alpha1.CreateServerRequestObject,
) (serverv1alpha1.CreateServerResponseObject, error) {
	if req.Body.Type == serverv1alpha1.Persistent {
		return nil, fmt.Errorf("not implemented")
	}

	args := server.CreateEphemeralServerArgs{
		Name:   req.Body.Name,
		Labels: ptr.Deref(req.Body.Labels),
	}

	serv, err := s.serverService.CreateEphemeralServer(ctx, args)
	if err != nil {
		return nil, fmt.Errorf("ehphemeral server: %w", err)
	}

	return serverv1alpha1.CreateServer202JSONResponse{
		JobCreateResponseJSONResponse: serverv1alpha1.JobCreateResponseJSONResponse{
			ResourceId: serv.ID,
			State:      serverv1alpha1.JobStateRunning, // TODO: pending?
		},
	}, nil
}

func (s Handler) DeleteServer(ctx context.Context, req serverv1alpha1.DeleteServerRequestObject) (serverv1alpha1.DeleteServerResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s Handler) GetServer(ctx context.Context, req serverv1alpha1.GetServerRequestObject) (serverv1alpha1.GetServerResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s Handler) UpdateServer(ctx context.Context, req serverv1alpha1.UpdateServerRequestObject) (serverv1alpha1.UpdateServerResponseObject, error) {
	//TODO implement me
	panic("implement me")
}
