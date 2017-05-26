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
	query :=
		`SELECT
			id as ID,
			name as Name,
			parent_id as ParentId
		FROM
			category`
	categories := []model.Category{}
	err := cdm.db.Select(&categories, query)
	if err != nil {
		return nil, err
	} else {
		return &categories, nil
	}
}

func (cdm *CategoryRepository) FindCategoryById(id int) (*model.Category, error) {
	query :=
		`SELECT
			id as ID,
			name as Name,
			parent_id as ParentId
		FROM
			category
		WHERE
			id=$1`
	categories := []model.Category{}
	err := cdm.db.Select(&categories, query, id)
	if err != nil {
		return nil, err
	} else {
		return &categories[0], nil
	}
}
