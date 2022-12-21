package rental

import (
	"context"
	"fmt"
	"lms/protogen/rental_service"
	"lms/storage"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RentalService is a struct that implements the server interface
type rentalService struct {
	stg storage.StorageI
	rental_service.UnimplementedRentalServiceServer
}

func NewRentalService(stg storage.StorageI) *rentalService {
	return &rentalService{
		stg: stg,
	}
}

// CreateRental ...
func (a *rentalService) CreateRental(c context.Context, req *rental_service.CreateRentalRequest) (*rental_service.Rental, error) {
	id := uuid.New()
	fmt.Println(id.String())
	err := a.stg.CreateRental(id.String(), req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.CreateRental: %s", err.Error())
	}
	res, err := a.stg.GetRentalById(id.String())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.Stg.GetRentalById: %s", err.Error())
	}
	return res, nil
}

// GetRentalById ...
func (a *rentalService) GetRentalById(c context.Context, req *rental_service.GetRentalByIdRequest) (*rental_service.Rental, error) {
	res, err := a.stg.GetRentalById(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "a.stg.GetRentalById: %s", err.Error())
	}
	return res, nil
}

// GetRentalList ...
func (a *rentalService) GetRentalList(c context.Context, req *rental_service.GetRentalListRequest) (*rental_service.GetRentalListResponse, error) {
	res, err := a.stg.GetRentalList(int(req.Offset), int(req.Limit), req.Search)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "a.stg.GetRentalList: %s", err.Error())
	}
	return res, nil
}

// UpdateRental ...
func (a *rentalService) UpdateRental(c context.Context, req *rental_service.UpdateRentalRequest) (*rental_service.Rental, error) {

	err := a.stg.UpdateRental(req)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "a.stg.UpdateRental: %s", err.Error())
	}

	res, e := a.stg.GetRentalById(req.Id)
	if e != nil {
		return nil, status.Errorf(codes.NotFound, "a.stg.UpdateRental: %s", e.Error())
	}
	return res, nil
}

// DeleteRental ...
func (a *rentalService) DeleteRental(c context.Context, req *rental_service.DeleteRentalRequest) (*rental_service.DeleteRentalResponse, error) {
	res, e := a.stg.GetRentalById(req.Id)
	if e != nil {
		return nil, status.Errorf(codes.NotFound, "a.stg.UpdateRental: %s", e.Error())
	}
	err := a.stg.DeleteRental(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "a.stg.DeleteRental: %s", err.Error())
	}

	return &rental_service.DeleteRentalResponse{
		RentalStatus: res.RentalStatus,
	}, nil
}
