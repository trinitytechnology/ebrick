package environment

import (
	"github.com/linkifysoft/ebrick/repository"
	"gorm.io/gorm"
)

type EnvironmentRepository interface {
	repository.CrudRepository[Environment]
}

type envOrgRepo struct {
	repository.CrudRepository[Environment]
	db *gorm.DB
}

func NewRepository(db *gorm.DB) EnvironmentRepository {
	return &envOrgRepo{
		CrudRepository: repository.NewCrudRepository[Environment](db),
		db:             db,
	}
}
