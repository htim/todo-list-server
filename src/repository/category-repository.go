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
	UpdateCategory(category model.Category) error
	DeleteCategory(id int) error
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
	var id int
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
	err = tx.QueryRow(query, category.Name, category.ParentId).Scan(&id)
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

func (cr *CategoryRepository) UpdateCategory(category model.Category) error {
	query :=
		`UPDATE
			category
		SET
			name=$1,
			parent_id=$2
		WHERE
			id=$3`
	tx, err := cr.db.Beginx()
	if err != nil {
		return err
	}
	_, err = tx.Exec(query, category.Name, category.ParentId, category.ID)
	if err != nil {
		multierror.Append(err, tx.Rollback())
		return err
	}
	tx.Commit()
	return nil
}

func (cr *CategoryRepository) DeleteCategory(id int) error {
	query :=
		`DELETE
		FROM
			category
		WHERE
			id=$1`
	tx, err := cr.db.Beginx()
	if err != nil {
		return err
	}
	_, err = tx.Exec(query, id)
	if err != nil {
		multierror.Append(err, tx.Rollback())
		return err
	}
	tx.Commit()
	return nil
}