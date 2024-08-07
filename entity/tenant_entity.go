package entity

import (
	"github.com/google/uuid"
)

type TenantAuditEntity struct {
	AuditEntity
	TenantId uuid.UUID `gorm:"type:uuid;index; not null" json:"tenant_id" validate:"required"`
}
