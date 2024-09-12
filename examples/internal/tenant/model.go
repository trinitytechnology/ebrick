package tenant

import "github.com/trinitytechnology/ebrick/entity"

type Tenant struct {
	entity.AuditEntity
	Name string `json:"name"`
}
