package service

import (
	"errors"

	"github.com/maftunhasansiam/SpotSync/dto"
	"github.com/maftunhasansiam/SpotSync/models"
	"github.com/maftunhasansiam/SpotSync/repository"
	"gorm.io/gorm"
)

type ReservationService interface {
	CreateReservation(userID uint, req *dto.CreateReservationRequest) (*dto.ReservationResponse, error)
	GetMyReservations(userID uint) ([]dto.MyReservationResponse, error)
	CancelReservation(reservationID uint, userID uint, role string) error
	GetAllReservations() ([]dto.AdminReservationResponse, error)
}
type reservationService struct {
	reservationRepo repository.ReservationRepository
	zoneRepo        repository.ZoneRepository
}
func NewReservationService(
	reservationRepo repository.ReservationRepository,
	zoneRepo repository.ZoneRepository,
) ReservationService {
	return &reservationService{
		reservationRepo: reservationRepo,
		zoneRepo:        zoneRepo,
	}
}

func (s *reservationService) CreateReservation(userID uint, req *dto.CreateReservationRequest) (*dto.ReservationResponse, error) {
	
	_, err := s.zoneRepo.FindByID(req.ZoneID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("parking zone not found")
		}
		return nil, errors.New("failed to retrieve parking zone")
	}
	reservation := &models.Reservation{
		UserID:       userID,
		ZoneID:       req.ZoneID,
		LicensePlate: req.LicensePlate,
		Status:       "active",
	}

	
	if err := s.reservationRepo.CreateWithLock(reservation); err != nil {
		if errors.Is(err, repository.ErrZoneFull) {
			return nil, repository.ErrZoneFull
		}
		return nil, errors.New("failed to create reservation")
	}

	return &dto.ReservationResponse{
		ID:           reservation.ID,
		UserID:       reservation.UserID,
		ZoneID:       reservation.ZoneID,
		LicensePlate: reservation.LicensePlate,
		Status:       reservation.Status,
		CreatedAt:    reservation.CreatedAt,
		UpdatedAt:    reservation.UpdatedAt,
	}, nil
}

func (s *reservationService) GetMyReservations(userID uint) ([]dto.MyReservationResponse, error) {
	reservations, err := s.reservationRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("failed to retrieve reservations")
	}

	var responses []dto.MyReservationResponse
	for _, r := range reservations {
		responses = append(responses, dto.MyReservationResponse{
			ID:           r.ID,
			LicensePlate: r.LicensePlate,
			Status:       r.Status,
			Zone: dto.ReservationZoneInfo{
				ID:   r.Zone.ID,
				Name: r.Zone.Name,
				Type: r.Zone.Type,
			},
			CreatedAt: r.CreatedAt,
		})
	}

	if responses == nil {
		responses = []dto.MyReservationResponse{}
	}

	return responses, nil
}

func (s *reservationService) CancelReservation(reservationID uint, userID uint, role string) error {
	reservation, err := s.reservationRepo.FindByID(reservationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("reservation not found")
		}
		return errors.New("failed to retrieve reservation")
	}
	if role != "admin" && reservation.UserID != userID {
		return errors.New("forbidden: you can only cancel your own reservations")
	}

	if reservation.Status != "active" {
		return errors.New("only active reservations can be cancelled")
	}

	if err := s.reservationRepo.Cancel(reservation); err != nil {
		return errors.New("failed to cancel reservation")
	}

	return nil
}

func (s *reservationService) GetAllReservations() ([]dto.AdminReservationResponse, error) {
	reservations, err := s.reservationRepo.FindAll()
	if err != nil {
		return nil, errors.New("failed to retrieve reservations")
	}

	var responses []dto.AdminReservationResponse
	for _, r := range reservations {
		responses = append(responses, dto.AdminReservationResponse{
			ID:           r.ID,
			LicensePlate: r.LicensePlate,
			Status:       r.Status,
			User: dto.ReservationUserInfo{
				ID:    r.User.ID,
				Name:  r.User.Name,
				Email: r.User.Email,
				Role:  r.User.Role,
			},
			Zone: dto.ReservationZoneInfo{
				ID:   r.Zone.ID,
				Name: r.Zone.Name,
				Type: r.Zone.Type,
			},
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		})
	}
	if responses == nil {
		responses = []dto.AdminReservationResponse{}
	}
	return responses, nil
}
