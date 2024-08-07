package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuditEntity struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	CreatedBy string         `json:"created_by"`
	UpdatedAt time.Time      `json:"updated_at"`
	UpdatedBy string         `json:"updated_by"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (am *AuditEntity) BeforeCreate(tx *gorm.DB) (err error) {
	if am.ID == uuid.Nil {
		am.ID = uuid.New()
	}
	return
}
