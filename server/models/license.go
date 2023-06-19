package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/mmdaz/lime/config"
)

// License is a ...
type License struct {
	ID             uint32    `gorm:"primary_key;auto_increment" json:"id"`
	SubscriptionID uint32    `sql:"type:int REFERENCES subscriptions(id)" json:"subscription_id"`
	License        []byte    `gorm:"null" json:"license"`
	Hash           string    `gorm:"null" json:"hash"`
	Status         bool      `gorm:"false" json:"status"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// SaveLicense is a ...
func (l *License) SaveLicense() (*License, error) {
	err := config.DB.Create(&l).Error
	if err != nil {
		return &License{}, err
	}
	return l, nil
}

// FindLicenseByID is a ...
func (l *License) FindLicenseByID(uid uint32) (*License, error) {
	err := config.DB.Model(License{}).Where("id = ?", uid).Take(&l).Error
	if err != nil {
		return &License{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &License{}, ErrKeyNotFound
	}
	return l, err
}

// FindLicense is a ...
func (l *License) FindLicense(key []byte) (*License, error) {
	err := config.DB.Model(License{}).Where("license = ?", key).Take(&l).Error
	if err != nil {
		return &License{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &License{}, ErrLicenseNotFound
	}
	return l, err
}

// UpdateLicense is a ...
func (l *License) UpdateLicense(uid uint32) (*License, error) {
	db := config.DB.Model(&License{}).Where("id = ?", uid).Take(&License{}).UpdateColumns(
		map[string]interface{}{
			"subscription_id": l.SubscriptionID,
			"license":         l.License,
			"status":          l.Status,
			"update_at":       time.Now(),
		},
	)
	if db.Error != nil {
		return &License{}, db.Error
	}

	err := config.DB.Model(&License{}).Where("id = ?", uid).Take(&l).Error
	if err != nil {
		return &License{}, err
	}
	return l, nil
}

// DeleteLicense is a ...
func (l *License) DeleteLicense(uid uint32) (int64, error) {
	db := config.DB.Model(&License{}).Where("id = ?", uid).Take(&License{}).Delete(&License{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

// DeactivateLicenseBySubID is a ...
func DeactivateLicenseBySubID(uid uint32) error {
	db := config.DB.Model(&License{}).Where("subscription_id = ?", uid).UpdateColumns(
		map[string]interface{}{
			"status":    false,
			"update_at": time.Now(),
		},
	)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

// LicensesListBySubscriptionID is a ...
func LicensesListBySubscriptionID(uid string) *[]License {
	licenses := []License{}
	db := config.DB.Where("subscription_id = ?", uid).Order("created_at DESC").Find(&licenses)
	if db.Error != nil {
		return &licenses
	}
	return &licenses
}
