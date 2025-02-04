package repositories

import (
	"gorm.io/gorm"
)

type GenericRepository[T any] struct {
	DB *gorm.DB
}

func NewGenericRepository[T any](db *gorm.DB) *GenericRepository[T] {
	return &GenericRepository[T]{
		DB: db,
	}
}

func (r *GenericRepository[T]) Create(entity T) error {
	return r.DB.Create(&entity).Error
}

func (r *GenericRepository[T]) FindByID(id uint) (*T, error) {
	var entity T
	err := r.DB.First(&entity, id).Error
	return &entity, err
}

func (r *GenericRepository[T]) FindAll() ([]T, error) {
	var entities []T
	err := r.DB.Find(&entities).Error
	return entities, err
}

func (r *GenericRepository[T]) Update(entity T) error {
	return r.DB.Save(&entity).Error
}

func (r *GenericRepository[T]) Delete(id uint) error {
	var entity T
	return r.DB.Delete(&entity, id).Error
}
