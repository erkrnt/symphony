package block

import (
	"context"

	"github.com/erkrnt/symphony/api"
	"github.com/erkrnt/symphony/internal/service"
	"github.com/erkrnt/symphony/internal/utils"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCServerControl : control service struct
type GRPCServerControl struct {
	Block *Block
}

// ServiceNew : handles creating a new service definition
func (s *GRPCServerControl) ServiceNew(ctx context.Context, in *api.RequestServiceNew) (*api.ResponseServiceNew, error) {
	serviceID := s.Block.Service.Key.ServiceID

	if serviceID != nil {
		st := status.New(codes.AlreadyExists, "service_already_exists")

		return nil, st.Err()
	}

	apiserverAddr := s.Block.Service.Flags.ManagerAddr

	if apiserverAddr == nil {
		st := status.New(codes.InvalidArgument, "invalid_apiserver_addr")

		return nil, st.Err()
	}

	conn, err := utils.NewClientConnTcp(apiserverAddr.String())

	if err != nil {
		st := status.New(codes.Internal, err.Error())

		return nil, st.Err()
	}

	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), utils.ContextTimeout)

	defer cancel()

	c := api.NewManagerClient(conn)

	serviceAddr := s.Block.Service.Flags.ServiceAddr.String()

	newServiceOptions := &api.RequestNewService{
		ServiceAddr: serviceAddr,
		ServiceType: api.ServiceType_BLOCK,
	}

	newService, err := c.NewService(ctx, newServiceOptions)

	if err != nil {
		st := status.New(codes.Internal, err.Error())

		return nil, st.Err()
	}

	newServiceID, err := uuid.Parse(newService.ID)

	if err != nil {
		st := status.New(codes.InvalidArgument, err.Error())

		return nil, st.Err()
	}

	saveKeyOptions := service.SaveKeyOptions{
		ConfigDir: s.Block.Service.Flags.ConfigDir,
		ServiceID: newServiceID,
	}

	saveErr := s.Block.Service.Key.Save(saveKeyOptions)

	if saveErr != nil {
		st := status.New(codes.Internal, saveErr.Error())

		return nil, st.Err()
	}

	res := &api.ResponseServiceNew{
		ServiceID: newService.ID,
	}

	return res, nil
}
