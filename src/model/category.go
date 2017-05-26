package model

type Category struct {
	ID int `db:"id"`
	Name string `db:"name"`
	ParentId int `db:"parent_id"`
}
