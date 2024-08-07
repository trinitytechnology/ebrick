package environment

import "github.com/linkifysoft/ebrick/entity"

type Environment struct {
	entity.TenantAuditEntity
	Name string `json:"name"`
}
