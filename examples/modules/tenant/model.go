package tenant

import "github.com/linkifysoft/ebrick/entity"

type Tenant struct {
	entity.AuditEntity
	Name string `json:"name"`
}
