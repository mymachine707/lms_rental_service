package storage

import "lms/protogen/rental_service"

type StorageI interface {
	CreateRental(id string, req *rental_service.CreateRentalRequest) error
	GetRentalById(id string) (*rental_service.Rental, error)
	GetRentalList(offset, limit int, search string) (*rental_service.GetRentalListResponse, error)
	UpdateRental(entity *rental_service.UpdateRentalRequest) error
	DeleteRental(id string) error
}
