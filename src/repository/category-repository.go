package repository

import (
	"model"
	"github.com/jmoiron/sqlx"
	"github.com/hashicorp/go-multierror"
)

type ICategoryRepository interface {
	FindAllCategories() (*[]model.Category, error)
	FindCategoryById(id int) (*model.Category, error)
	CreateCategory(category model.Category) (int, error)
}

type CategoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (cr *CategoryRepository) FindAllCategories() (*[]model.Category, error) {
	query :=
		`SELECT
			id as ID,
			name as Name,
			parent_id as ParentId
		FROM
			category`
	categories := []model.Category{}
	err := cr.db.Select(&categories, query)
	if err != nil {
		return nil, err
	} else {
		return &categories, nil
	}
}

func (cr *CategoryRepository) FindCategoryById(id int) (*model.Category, error) {
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
	err := cr.db.Select(&categories, query, id)
	if err != nil {
		return nil, err
	} else {
		return &categories[0], nil
	}
}

func (cr *CategoryRepository) CreateCategory(category model.Category) (int, error) {
	query :=
		`INSERT INTO
			category (name, parent_id)
		VALUES
			($1, $2)
		RETURNING id`
	tx, err := cr.db.Beginx()
	if err != nil {
		return -1, err
	}
	res, err := tx.Exec(query, category.Name, category.ParentId)
	if err != nil {
		multierror.Append(err, tx.Rollback())
		return -1, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		multierror.Append(err, tx.Rollback())
		return -1, err
	}
	err = tx.Commit()
	if err != nil {
		return -1, err
	}
	return int(id), err
}
