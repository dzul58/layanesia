package repositories

import (
	"layanesia-server/internal/models"

	"gorm.io/gorm"
)

type LocationRepository interface {
	GetProvinces() ([]string, error)
	GetCitiesByProvince(province string) ([]string, error)
	GetDistrictsByCity(city string) ([]models.Location, error)
	SearchLocations(query string) ([]models.Location, error)
}

type locationRepository struct {
	db *gorm.DB
}

func NewLocationRepository(db *gorm.DB) LocationRepository {
	return &locationRepository{db: db}
}

func (r *locationRepository) GetProvinces() ([]string, error) {
	var provinces []string
	err := r.db.Model(&models.Location{}).
		Distinct("province").
		Order("province ASC").
		Pluck("province", &provinces).Error
	return provinces, err
}

func (r *locationRepository) GetCitiesByProvince(province string) ([]string, error) {
	var cities []string
	err := r.db.Model(&models.Location{}).
		Where("LOWER(province) = LOWER(?)", province).
		Distinct("city_or_regency").
		Order("city_or_regency ASC").
		Pluck("city_or_regency", &cities).Error
	return cities, err
}

func (r *locationRepository) GetDistrictsByCity(city string) ([]models.Location, error) {
	var locations []models.Location
	err := r.db.Where("LOWER(city_or_regency) = LOWER(?)", city).
		Order("district ASC").
		Find(&locations).Error
	return locations, err
}

func (r *locationRepository) SearchLocations(query string) ([]models.Location, error) {
	var locations []models.Location
	searchPattern := "%" + query + "%"
	err := r.db.Where("LOWER(province) LIKE LOWER(?) OR LOWER(city_or_regency) LIKE LOWER(?) OR LOWER(district) LIKE LOWER(?)",
		searchPattern, searchPattern, searchPattern).
		Limit(20).
		Find(&locations).Error
	return locations, err
}
