package apiserver

import (
	"context"
	"fmt"
	"net"

	"github.com/erkrnt/symphony/api"
	"github.com/erkrnt/symphony/internal/service"
	"github.com/google/uuid"
	consul "github.com/hashicorp/consul/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCServerAPIServer : APIServer GRPC endpoints
type GRPCServerAPIServer struct {
	APIServer *APIServer
}

// GetLogicalVolume : retrieves a logical volume from state
func (s *GRPCServerAPIServer) GetLogicalVolume(ctx context.Context, in *api.RequestLogicalVolume) (*api.LogicalVolume, error) {
	lvID, err := uuid.Parse(in.ID)

	if err != nil {
		st := status.New(codes.InvalidArgument, err.Error())

		return nil, st.Err()
	}

	lv, err := s.APIServer.getLogicalVolumeByID(lvID)

	if err != nil {
		st := status.New(codes.Internal, err.Error())

		return nil, st.Err()
	}

	if lv == nil {
		st := status.New(codes.NotFound, "invalid_logical_volume_id")

		return nil, st.Err()
	}

	return lv, nil
}

// GetLogicalVolumes : retrieves all logical volumes from state
func (s *GRPCServerAPIServer) GetLogicalVolumes(ctx context.Context, in *api.RequestLogicalVolumes) (*api.ResponseLogicalVolumes, error) {
	lvs, err := s.APIServer.getLogicalVolumes()

	if err != nil {
		st := status.New(codes.Internal, err.Error())

		return nil, st.Err()
	}

	res := &api.ResponseLogicalVolumes{
		Results: lvs,
	}

	return res, nil
}

// GetPhysicalVolume : retrieves a physical volume from state
func (s *GRPCServerAPIServer) GetPhysicalVolume(ctx context.Context, in *api.RequestPhysicalVolume) (*api.PhysicalVolume, error) {
	pvID, err := uuid.Parse(in.ID)

	if err != nil {
		st := status.New(codes.InvalidArgument, err.Error())

		return nil, st.Err()
	}

	pv, err := s.APIServer.getPhysicalVolumeByID(pvID)

	if err != nil {
		st := status.New(codes.Internal, err.Error())

		return nil, st.Err()
	}

	if pv == nil {
		st := status.New(codes.NotFound, "invalid_physical_volume_id")

		return nil, st.Err()
	}

	return pv, nil
}

// GetPhysicalVolumes : retrieves all physical volumes from state
func (s *GRPCServerAPIServer) GetPhysicalVolumes(ctx context.Context, in *api.RequestPhysicalVolumes) (*api.ResponsePhysicalVolumes, error) {
	pvs, err := s.APIServer.getPhysicalVolumes()

	if err != nil {
		st := status.New(codes.Internal, err.Error())
		return nil, st.Err()
	}

	res := &api.ResponsePhysicalVolumes{
		Results: pvs,
	}

	return res, nil
}

// GetService : retrieves a service from state
func (s *GRPCServerAPIServer) GetService(ctx context.Context, in *api.RequestService) (*api.Service, error) {
	serviceID, err := uuid.Parse(in.ServiceID)

	if err != nil {
		st := status.New(codes.InvalidArgument, "invalid_service_id")

		return nil, st.Err()
	}

	service, err := s.APIServer.getServiceByID(serviceID)

	if err != nil {
		st := status.New(codes.Internal, err.Error())

		return nil, st.Err()
	}

	if service == nil {
		st := status.New(codes.NotFound, "invalid_service")

		return nil, st.Err()
	}

	return service, nil
}

// GetServices : retrieves all services from state
func (s *GRPCServerAPIServer) GetServices(ctx context.Context, in *api.RequestServices) (*api.ResponseServices, error) {
	services, err := s.APIServer.getServices()

	if err != nil {
		st := status.New(codes.Internal, err.Error())

		return nil, st.Err()
	}

	res := &api.ResponseServices{
		Results: services,
	}

	return res, nil
}

// GetVolumeGroup : retrieves a volume group from state
func (s *GRPCServerAPIServer) GetVolumeGroup(ctx context.Context, in *api.RequestVolumeGroup) (*api.VolumeGroup, error) {
	vgID, err := uuid.Parse(in.ID)

	if err != nil {
		st := status.New(codes.InvalidArgument, err.Error())
		return nil, st.Err()
	}

	vg, err := s.APIServer.getVolumeGroupByID(vgID)

	if err != nil {
		st := status.New(codes.Internal, err.Error())
		return nil, st.Err()
	}

	if vg == nil {
		st := status.New(codes.NotFound, "invalid_volume_group_id")
		return nil, st.Err()
	}

	return vg, nil
}

// GetVolumeGroups : retrieves all volume groups from state
func (s *GRPCServerAPIServer) GetVolumeGroups(ctx context.Context, in *api.RequestVolumeGroups) (*api.ResponseVolumeGroups, error) {
	vgs, err := s.APIServer.getVolumeGroups()

	if err != nil {
		st := status.New(codes.Internal, err.Error())
		return nil, st.Err()
	}

	res := &api.ResponseVolumeGroups{
		Results: vgs,
	}

	return res, nil
}

// NewLogicalVolume : creates a new logical volume in state
func (s *GRPCServerAPIServer) NewLogicalVolume(ctx context.Context, in *api.RequestNewLogicalVolume) (*api.LogicalVolume, error) {
	vgID, err := uuid.Parse(in.VolumeGroupID)

	if err != nil {
		st := status.New(codes.InvalidArgument, err.Error())

		return nil, st.Err()
	}

	vg, err := s.APIServer.getVolumeGroupByID(vgID)

	if err != nil {
		st := status.New(codes.Internal, err.Error())

		return nil, st.Err()
	}

	if vg == nil {
		st := status.New(codes.NotFound, "invalid_volume_group_id")

		return nil, st.Err()
	}

	pvID, err := uuid.Parse(vg.PhysicalVolumeID)

	if err != nil {
		st := status.New(codes.InvalidArgument, "invalid_service_id")

		return nil, st.Err()
	}

	pv, err := s.APIServer.getPhysicalVolumeByID(pvID)

	if err != nil {
		st := status.New(codes.Internal, err.Error())

		return nil, st.Err()
	}

	if pv == nil {
		st := status.New(codes.NotFound, "invalid_physical_volume_id")

		return nil, st.Err()
	}

	lvID := uuid.New()

	lv := &api.LogicalVolume{
		ID:            lvID.String(),
		Size:          in.Size,
		Status:        api.ResourceStatus_REVIEW_IN_PROGRESS,
		VolumeGroupID: vg.ID,
	}

	saveErr := s.APIServer.saveLogicalVolume(lv)

	if saveErr != nil {
		return nil, saveErr
	}

	return lv, nil
}

// NewPhysicalVolume : creates a new physical volume in state
func (s *GRPCServerAPIServer) NewPhysicalVolume(ctx context.Context, in *api.RequestNewPhysicalVolume) (*api.PhysicalVolume, error) {
	serviceID, err := uuid.Parse(in.ServiceID)

	if err != nil {
		st := status.New(codes.InvalidArgument, err.Error())
		return nil, st.Err()
	}

	service, err := s.APIServer.getServiceByID(serviceID)

	if service == nil {
		st := status.New(codes.NotFound, "invalid_service_id")
		return nil, st.Err()
	}

	volumes, err := s.APIServer.getPhysicalVolumes()

	if err != nil {
		st := status.New(codes.Internal, err.Error())
		return nil, st.Err()
	}

	var volume *api.PhysicalVolume

	for _, v := range volumes {
		if in.DeviceName == v.DeviceName && service.ID == v.ServiceID {
			volume = v
		}
	}

	if volume != nil {
		st := status.New(codes.AlreadyExists, "physical_volume_exists")
		return nil, st.Err()
	}

	pvID := uuid.New()

	pv := &api.PhysicalVolume{
		DeviceName: in.DeviceName,
		ID:         pvID.String(),
		ServiceID:  service.ID,
	}

	saveErr := s.APIServer.savePhysicalVolume(pv)

	if saveErr != nil {
		return nil, saveErr
	}

	return pv, nil
}

// NewService : creates a new service in state
func (s *GRPCServerAPIServer) NewService(ctx context.Context, in *api.RequestNewService) (*api.Service, error) {
	client, err := service.NewConsulClient(s.APIServer.Flags.ConsulAddr)

	if err != nil {
		st := status.New(codes.Internal, err.Error())

		return nil, st.Err()
	}

	agent := client.Agent()

	regName := uuid.New()

	regCheckGRPC, err := net.ResolveTCPAddr("tcp", in.HealthServiceAddr)

	if err != nil {
		st := status.New(codes.InvalidArgument, err.Error())

		return nil, st.Err()
	}

	regCheck := &consul.AgentServiceCheck{
		CheckID:    regName.String(),
		GRPC:       regCheckGRPC.String(),
		GRPCUseTLS: false,
		Interval:   "10s",
		Name:       regName.String(),
	}

	reg := &consul.AgentServiceRegistration{
		Address: regCheckGRPC.IP.String(),
		Check:   regCheck,
		Name:    regName.String(),
		Tags:    []string{in.ServiceType.String()},
	}

	regErr := agent.ServiceRegister(reg)

	if regErr != nil {
		st := status.New(codes.Internal, regErr.Error())

		return nil, st.Err()
	}

	serviceQueryOpts := &consul.QueryOptions{}

	regService, _, err := agent.Service(reg.Name, serviceQueryOpts)

	if err != nil {
		st := status.New(codes.Internal, err.Error())

		return nil, st.Err()
	}

	service := &api.Service{
		ID:   regService.ID,
		Type: in.ServiceType,
	}

	// TODO: open stream to apiserver for events

	return service, nil
}

// NewVolumeGroup : creates a new volume group in state
func (s *GRPCServerAPIServer) NewVolumeGroup(ctx context.Context, in *api.RequestNewVolumeGroup) (*api.VolumeGroup, error) {
	physicalVolumeID, err := uuid.Parse(in.PhysicalVolumeID)

	if err != nil {
		st := status.New(codes.InvalidArgument, err.Error())
		return nil, st.Err()
	}

	physicalVolume, err := s.APIServer.getPhysicalVolumeByID(physicalVolumeID)

	if err != nil {
		st := status.New(codes.Internal, err.Error())
		return nil, st.Err()
	}

	if physicalVolume == nil {
		st := status.New(codes.NotFound, "invalid_physical_volume_id")
		return nil, st.Err()
	}

	physicalVolumeServiceID, err := uuid.Parse(physicalVolume.ID)

	if err != nil {
		st := status.New(codes.InvalidArgument, "invalid_service_id")
		return nil, st.Err()
	}

	service, err := s.APIServer.getServiceByID(physicalVolumeServiceID)

	if service == nil {
		st := status.New(codes.NotFound, "invalid_service_id")
		return nil, st.Err()
	}

	volumeGroupID := uuid.New()

	vg := &api.VolumeGroup{
		ID:               volumeGroupID.String(),
		PhysicalVolumeID: physicalVolume.ID,
		Status:           api.ResourceStatus_REVIEW_IN_PROGRESS,
	}

	saveErr := s.APIServer.saveVolumeGroup(vg)

	if saveErr != nil {
		return nil, saveErr
	}

	return vg, nil
}

// RemoveLogicalVolume : removes a logical volume from state
func (s *GRPCServerAPIServer) RemoveLogicalVolume(ctx context.Context, in *api.RequestLogicalVolume) (*api.ResponseStatus, error) {
	logicalVolumeID, err := uuid.Parse(in.ID)

	if err != nil {
		st := status.New(codes.InvalidArgument, err.Error())
		return nil, st.Err()
	}

	lv, err := s.APIServer.getLogicalVolumeByID(logicalVolumeID)

	if err != nil {
		st := status.New(codes.Internal, err.Error())
		return nil, st.Err()
	}

	if lv == nil {
		st := status.New(codes.NotFound, "invalid_logical_volume_id")
		return nil, st.Err()
	}

	resourceKey := fmt.Sprintf("/LogicalVolume/%s", lv.ID)

	delErr := s.APIServer.removeResource(resourceKey)

	if delErr != nil {
		st := status.New(codes.Internal, delErr.Error())

		return nil, st.Err()
	}

	// TODO: emit an event change to block services

	res := &api.ResponseStatus{SUCCESS: true}

	return res, nil
}

// RemovePhysicalVolume : removes a physical volume from state
func (s *GRPCServerAPIServer) RemovePhysicalVolume(ctx context.Context, in *api.RequestPhysicalVolume) (*api.ResponseStatus, error) {
	physicalVolumeID, err := uuid.Parse(in.ID)

	if err != nil {
		st := status.New(codes.InvalidArgument, err.Error())

		return nil, st.Err()
	}

	pv, err := s.APIServer.getPhysicalVolumeByID(physicalVolumeID)

	if err != nil {
		st := status.New(codes.Internal, err.Error())

		return nil, st.Err()
	}

	if pv == nil {
		st := status.New(codes.NotFound, "invalid_physical_volume_id")

		return nil, st.Err()
	}

	serviceID, err := uuid.Parse(pv.ServiceID)

	if err != nil {
		st := status.New(codes.InvalidArgument, err.Error())

		return nil, st.Err()
	}

	service, err := s.APIServer.getServiceByID(serviceID)

	if err != nil {
		st := status.New(codes.Internal, err.Error())

		return nil, st.Err()
	}

	if service == nil {
		st := status.New(codes.NotFound, "invalid_service_id")

		return nil, st.Err()
	}

	resourceKey := fmt.Sprintf("/PhysicalVolume/%s", pv.ID)

	delRes := s.APIServer.removeResource(resourceKey)

	if delRes != nil {
		st := status.New(codes.Internal, delRes.Error())

		return nil, st.Err()
	}

	res := &api.ResponseStatus{SUCCESS: true}

	return res, nil
}

// RemoveService : removes a service from state
func (s *GRPCServerAPIServer) RemoveService(ctx context.Context, in *api.RequestService) (*api.ResponseStatus, error) {
	serviceID, err := uuid.Parse(in.ServiceID)

	if err != nil {
		st := status.New(codes.InvalidArgument, err.Error())

		return nil, st.Err()
	}

	srvc, err := s.APIServer.getServiceByID(serviceID)

	if err != nil {
		st := status.New(codes.Internal, err.Error())

		return nil, st.Err()
	}

	resourceKey := fmt.Sprintf("/Service/%s", srvc.ID)

	delErr := s.APIServer.removeResource(resourceKey)

	if delErr != nil {
		st := status.New(codes.Internal, delErr.Error())

		return nil, st.Err()
	}

	res := &api.ResponseStatus{SUCCESS: true}

	return res, nil
}

// RemoveVolumeGroup : removes a volume group from state
func (s *GRPCServerAPIServer) RemoveVolumeGroup(ctx context.Context, in *api.RequestVolumeGroup) (*api.ResponseStatus, error) {
	volumeGroupID, err := uuid.Parse(in.ID)

	if err != nil {
		st := status.New(codes.InvalidArgument, err.Error())

		return nil, st.Err()
	}

	vg, err := s.APIServer.getVolumeGroupByID(volumeGroupID)

	if err != nil {
		st := status.New(codes.Internal, err.Error())

		return nil, st.Err()
	}

	if vg == nil {
		st := status.New(codes.NotFound, "invalid_volume_group_id")

		return nil, st.Err()
	}

	resourceKey := fmt.Sprintf("/VolumeGroup/%s", vg.ID)

	delErr := s.APIServer.removeResource(resourceKey)

	if delErr != nil {
		st := status.New(codes.Internal, delErr.Error())

		return nil, st.Err()
	}

	res := &api.ResponseStatus{SUCCESS: true}

	return res, nil
}