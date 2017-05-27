package model

type Category struct {
	ID int `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	ParentId int `json:"parentId,omitempty"`
}

