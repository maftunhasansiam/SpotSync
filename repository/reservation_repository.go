package repository

import (
	"errors"
	"github.com/maftunhasansiam/SpotSync/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrZoneFull = errors.New("parking zone is full")

type ReservationRepository interface {
	CreateWithLock(reservation *models.Reservation) error
	FindByUserID(userID uint) ([]models.Reservation, error)
	FindByID(id uint) (*models.Reservation, error)
	FindAll() ([]models.Reservation, error)
	Cancel(reservation *models.Reservation) error
}

type reservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) ReservationRepository {
	return &reservationRepository{db: db}
}

func (r *reservationRepository) CreateWithLock(reservation *models.Reservation) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var zone models.ParkingZone
		
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&zone, reservation.ZoneID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("parking zone not found")
			}
			return err
		}

		var activeCount int64
		if err := tx.Model(&models.Reservation{}).
			Where("zone_id = ? AND status = ?", reservation.ZoneID, "active").
			Count(&activeCount).Error; err != nil {
			return err
		}
		
		if activeCount >= int64(zone.TotalCapacity) {
			return ErrZoneFull
		}

		return tx.Create(reservation).Error
	})
}

func (r *reservationRepository) FindByUserID(userID uint) ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := r.db.Preload("Zone").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&reservations).Error
	return reservations, err
}

func (r *reservationRepository) FindByID(id uint) (*models.Reservation, error) {
	var reservation models.Reservation
	err := r.db.First(&reservation, id).Error
	if err != nil {
		return nil, err
	}
	return &reservation, nil
}

func (r *reservationRepository) FindAll() ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := r.db.Preload("User").Preload("Zone").
		Order("created_at DESC").
		Find(&reservations).Error
	return reservations, err
}
func (r *reservationRepository) Cancel(reservation *models.Reservation) error {
	reservation.Status = "cancelled"
	return r.db.Save(reservation).Error
}
