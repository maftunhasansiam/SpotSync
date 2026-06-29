package service

import (
	"errors"
	"github.com/maftunhasansiam/SpotSync/dto"
	"github.com/maftunhasansiam/SpotSync/models"
	"github.com/maftunhasansiam/SpotSync/repository"
	"gorm.io/gorm"
)

type ZoneService interface {
	CreateZone(req *dto.CreateZoneRequest) (*dto.ZoneResponse, error)
	GetAllZones() ([]dto.ZoneResponse, error)
	GetZoneByID(id uint) (*dto.ZoneResponse, error)
	UpdateZone(id uint, req *dto.UpdateZoneRequest) (*dto.ZoneResponse, error)
	DeleteZone(id uint) error
}

type zoneService struct {
	zoneRepo repository.ZoneRepository
}

func NewZoneService(zoneRepo repository.ZoneRepository) ZoneService {
	return &zoneService{zoneRepo: zoneRepo}
}

func (s *zoneService) CreateZone(req *dto.CreateZoneRequest) (*dto.ZoneResponse, error) {
	zone := &models.ParkingZone{
		Name:          req.Name,
		Type:          req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}

	if err := s.zoneRepo.Create(zone); err != nil {
		return nil, errors.New("failed to create parking zone")
	}

	return s.buildZoneResponse(zone, 0), nil
}

func (s *zoneService) GetAllZones() ([]dto.ZoneResponse, error) {
	zones, err := s.zoneRepo.FindAll()
	if err != nil {
		return nil, errors.New("failed to retrieve parking zones")
	}

	var responses []dto.ZoneResponse
	for _, zone := range zones {
		
		activeCount, err := s.zoneRepo.CountActiveReservations(zone.ID)
		if err != nil {
			activeCount = 0
		}
		availableSpots := zone.TotalCapacity - int(activeCount)
		if availableSpots < 0 {
			availableSpots = 0
		}
		responses = append(responses, *s.buildZoneResponse(&zone, availableSpots))
	}

	
	if responses == nil {
		responses = []dto.ZoneResponse{}
	}

	return responses, nil
}

func (s *zoneService) GetZoneByID(id uint) (*dto.ZoneResponse, error) {
	zone, err := s.zoneRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("parking zone not found")
		}
		return nil, errors.New("failed to retrieve parking zone")
	}

	activeCount, err := s.zoneRepo.CountActiveReservations(zone.ID)
	if err != nil {
		activeCount = 0
	}
	availableSpots := zone.TotalCapacity - int(activeCount)
	if availableSpots < 0 {
		availableSpots = 0
	}

	return s.buildZoneResponse(zone, availableSpots), nil
}

func (s *zoneService) UpdateZone(id uint, req *dto.UpdateZoneRequest) (*dto.ZoneResponse, error) {
	zone, err := s.zoneRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("parking zone not found")
		}
		return nil, errors.New("failed to retrieve parking zone")
	}


	if req.Name != "" {
		zone.Name = req.Name
	}
	if req.Type != "" {
		zone.Type = req.Type
	}
	if req.TotalCapacity > 0 {
		zone.TotalCapacity = req.TotalCapacity
	}
	if req.PricePerHour > 0 {
		zone.PricePerHour = req.PricePerHour
	}

	if err := s.zoneRepo.Update(zone); err != nil {
		return nil, errors.New("failed to update parking zone")
	}
	activeCount, _ := s.zoneRepo.CountActiveReservations(zone.ID)
	availableSpots := zone.TotalCapacity - int(activeCount)
	if availableSpots < 0 {
		availableSpots = 0
	}

	return s.buildZoneResponse(zone, availableSpots), nil
}

func (s *zoneService) DeleteZone(id uint) error {
	_, err := s.zoneRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("parking zone not found")
		}
		return errors.New("failed to retrieve parking zone")
	}
	if err := s.zoneRepo.Delete(id); err != nil {
		return errors.New("failed to delete parking zone")
	}
	return nil
}

func (s *zoneService) buildZoneResponse(zone *models.ParkingZone, availableSpots int) *dto.ZoneResponse {
	return &dto.ZoneResponse{
		ID:             zone.ID,
		Name:           zone.Name,
		Type:           zone.Type,
		TotalCapacity:  zone.TotalCapacity,
		AvailableSpots: availableSpots,
		PricePerHour:   zone.PricePerHour,
		CreatedAt:      zone.CreatedAt,
	}
}
