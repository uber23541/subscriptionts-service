package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Subscription struct {
	gorm.Model  `swaggerignore:"true"`
	ServiceName string     `gorm:"not null;uniqueIndex:uniq_user_service" json:"service_name"`
	Price       int        `gorm:"not null" json:"price" `
	UserID      uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:uniq_user_service" json:"user_id"`
	StartDate   time.Time  `gorm:"type:date;not null" json:"start_date"`
	EndDate     *time.Time `gorm:"type:date" json:"end_date,omitempty"`
}
