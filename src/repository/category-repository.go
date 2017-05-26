package repository

import (
	"model"
	"github.com/jmoiron/sqlx"
)

type ICategoryRepository interface {
	FindAllCategories() (*[]model.Category, error)
	FindCategoryById(id int) (*model.Category, error)
}

type CategoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (cdm *CategoryRepository) FindAllCategories() (*[]model.Category, error) {
	categories := []model.Category{}
	err := cdm.db.Select(&categories, "SELECT * FROM category")
	if err != nil {
		return nil, err
	} else {
		return &categories, nil
	}
}

func (cdm *CategoryRepository) FindCategoryById(id int) (*model.Category, error) {
	categories := []model.Category{}
	err := cdm.db.Select(&categories, "SELECT * FROM category WHERE id=$1", id)
	if err != nil {
		return nil, err
	} else {
		return &categories[0], nil
	}
}

