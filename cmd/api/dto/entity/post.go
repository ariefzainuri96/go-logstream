package entity

import (
	_ "gorm.io/gorm"
)

// @Model
type Post struct {
	BaseEntity
	ProjectId uint    `gorm:"type:int;not null;column:project_id" json:"project_id"`
	Project   Project `json:"project"`
	Title     string  `gorm:"type:varchar(255);not null;column:title" json:"title"`
	Content   string  `gorm:"type:text;not null;column:content" json:"content"`
	Category  string  `gorm:"type:varchar(50);not null;column:category" json:"category"` // 'feature', 'bugfix', 'maintenance'
	Status    string  `gorm:"type:varchar(20);not null;column:status" json:"status"` // 'draft', 'published'
}

/*
	for filtering field use like this for [carts] table:
	- carts.quantity -> even for current table filtering, always call the table name like this
	- products.name -> filter using products table with field name ->
	remember to not using struct field -> always use real tables and field name
*/

func (Post) TableName() string {
	return "posts"
}
