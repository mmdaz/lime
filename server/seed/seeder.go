package seed

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/mmdaz/lime/server/models"
)

var customers = []models.Customer{
	{
		ID:     1,
		Name:   "Alex Past",
		Status: true,
	},
	{
		ID:     2,
		Name:   "Lisa Boston",
		Status: true,
	},
	{
		ID:     3,
		Name:   "Adam Potar",
		Status: true,
	},
	{
		ID:     4,
		Name:   "Greg Gordon",
		Status: true,
	},
}

var subscription = []models.Subscription{
	{
		CustomerID: 1,
		Status:     true,
	},
	{
		CustomerID: 2,
		Status:     true,
	},
	{
		CustomerID: 3,
		Status:     true,
	},
	{
		CustomerID: 4,
		Status:     false,
	},
}

// Load import test data to database
func Load(db *gorm.DB) {
	err := db.DropTableIfExists(&models.Customer{}, &models.Subscription{}, &models.License{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.AutoMigrate(&models.Customer{}, &models.Subscription{}, &models.License{}).Error

	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	for i := range customers {
		err = db.Model(&models.Customer{}).Create(&customers[i]).Error
		if err != nil {
			log.Fatalf("cannot seed customer table: %v", err)
		}
	}

	for i := range subscription {
		err = db.Model(&models.Subscription{}).Create(&subscription[i]).Error
		if err != nil {
			log.Fatalf("cannot seed subscription table: %v", err)
		}
	}
}
