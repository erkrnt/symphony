package main

import (
	"context"

	"github.com/erkrnt/symphony/protobuff"
	"github.com/erkrnt/symphony/schema"
	"github.com/erkrnt/symphony/service"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes"
)

func (server *managerServer) GetServiceByHostname(ctx context.Context, in *protobuff.GetServiceByHostnameFields) (*protobuff.Service, error) {
	svc, err := getServiceByHostname(server.db, in.Hostname)

	if err != nil {
		return nil, service.HandleProtoError(err)
	}

	if svc == nil {
		err := status.Error(codes.NotFound, "invalid_service")
		return nil, service.HandleProtoError(err)
	}

	createdAt, err := ptypes.TimestampProto(svc.CreatedAt)

	if err != nil {
		return nil, service.HandleProtoError(err)
	}

	updatedAt, err := ptypes.TimestampProto(svc.UpdatedAt)

	if err != nil {
		return nil, service.HandleProtoError(err)
	}

	s := &protobuff.Service{
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
		ID:            svc.ID.String(),
		Hostname:      svc.Hostname,
		ServiceTypeID: svc.ServiceTypeID.String(),
	}

	logrus.WithFields(logrus.Fields{"Hostname": s.Hostname}).Info("GetServiceByHostname")

	return s, nil
}

func (server *managerServer) GetServiceByID(ctx context.Context, in *protobuff.GetServiceByIDFields) (*protobuff.Service, error) {
	id, err := uuid.Parse(in.ID)

	if err != nil {
		return nil, service.HandleProtoError(err)
	}

	svc, err := getServiceByID(server.db, id)

	if err != nil {
		return nil, service.HandleProtoError(err)
	}

	if svc == nil {
		err := status.Error(codes.NotFound, "invalid_service")
		return nil, service.HandleProtoError(err)
	}

	createdAt, err := ptypes.TimestampProto(svc.CreatedAt)

	if err != nil {
		return nil, service.HandleProtoError(err)
	}

	updatedAt, err := ptypes.TimestampProto(svc.UpdatedAt)

	if err != nil {
		return nil, service.HandleProtoError(err)
	}

	s := &protobuff.Service{
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
		ID:            svc.ID.String(),
		Hostname:      svc.Hostname,
		ServiceTypeID: svc.ServiceTypeID.String(),
	}

	logrus.WithFields(logrus.Fields{"ID": s.ID}).Info("GetServiceByID")

	return s, nil
}

func (server *managerServer) GetServiceTypeByName(ctx context.Context, in *protobuff.GetServiceTypeByNameFields) (*protobuff.ServiceType, error) {
	servicetype, err := getServiceTypeByName(server.db, in.Name)

	if err != nil {
		return nil, service.HandleProtoError(err)
	}

	if servicetype == nil {
		err := status.Error(codes.NotFound, "invalid_service_type")
		return nil, service.HandleProtoError(err)
	}

	st := &protobuff.ServiceType{
		ID:   servicetype.ID.String(),
		Name: servicetype.Name,
	}

	logrus.WithFields(logrus.Fields{"Name": st.Name}).Info("GetServiceTypeByName")

	return st, nil
}

func (server *managerServer) NewService(ctx context.Context, in *protobuff.NewServiceFields) (*protobuff.Service, error) {
	serviceTypeID, err := uuid.Parse(in.ServiceTypeID)

	if err != nil {
		return nil, service.HandleProtoError(err)
	}

	serviceType, err := getServiceTypeByID(server.db, serviceTypeID)

	if err != nil {
		return nil, service.HandleProtoError(err)
	}

	if serviceType == nil {
		err := status.Error(codes.NotFound, "invalid_service_type")
		return nil, service.HandleProtoError(err)
	}

	svc := schema.Service{Hostname: in.Hostname, ServiceTypeID: serviceType.ID}

	if err := server.db.Create(&svc).Error; err != nil {
		return nil, service.HandleProtoError(err)
	}

	createdAt, err := ptypes.TimestampProto(svc.CreatedAt)

	if err != nil {
		return nil, service.HandleProtoError(err)
	}

	updatedAt, err := ptypes.TimestampProto(svc.UpdatedAt)

	if err != nil {
		return nil, service.HandleProtoError(err)
	}

	s := &protobuff.Service{
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
		ID:            svc.ID.String(),
		Hostname:      svc.Hostname,
		ServiceTypeID: svc.ServiceTypeID.String(),
	}

	logFields := logrus.Fields{"Hostname": s.Hostname, "ServiceTypeID": s.ServiceTypeID}

	logrus.WithFields(logFields).Info("NewService")

	return s, nil
}