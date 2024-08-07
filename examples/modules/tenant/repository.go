package tenant

import (
	"github.com/linkifysoft/ebrick/repository"
	"gorm.io/gorm"
)

type TenantRepository interface {
	repository.CrudRepository[Tenant]
}

type tenantOrgRepo struct {
	repository.CrudRepository[Tenant]
	db *gorm.DB
}

func NewRepository(db *gorm.DB) TenantRepository {
	return &tenantOrgRepo{
		CrudRepository: repository.NewCrudRepository[Tenant](db),
		db:             db,
	}
}
