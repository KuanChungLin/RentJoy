package repositories

import (
	interfaces "rentjoy/internal/interfaces/repositories"
	"rentjoy/internal/models"

	"gorm.io/gorm"
)

type BudgetRepository struct {
	*GenericRepository[models.Budget]
}

func NewBudgetRepository(db *gorm.DB) interfaces.BudgetRepository {
	return &BudgetRepository{
		GenericRepository: NewGenericRepository[models.Budget](db),
	}
}
