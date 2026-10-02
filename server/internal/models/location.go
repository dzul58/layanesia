package models

type Location struct {
	ID            uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Country       string `gorm:"type:varchar(100);default:'Indonesia';not null" json:"country"`
	Province      string `gorm:"type:varchar(100);not null;index" json:"province"`
	CityOrRegency string `gorm:"type:varchar(100);not null;index" json:"city_or_regency"`
	District      string `gorm:"type:varchar(100);not null;index" json:"district"`
	PostalCode    *string `gorm:"type:varchar(20)" json:"postal_code,omitempty"`
}
