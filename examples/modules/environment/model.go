package environment

import "github.com/trinitytechnology/ebrick/entity"

type Environment struct {
	entity.TenantAuditEntity
	Name string `json:"name"`
}
