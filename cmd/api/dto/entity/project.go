package entity

import (
	_ "gorm.io/gorm"
)

// @Model
type Project struct {
	UpdateEntity
	UserId uint `gorm:"type:int;not null;column:user_id" json:"user_id"`
	// User       User   `json:"user"`
	Name            string `gorm:"type:varchar(255);not null;column:name" json:"name"`
	Slug            string `gorm:"type:varchar(255);not null;column:slug" json:"slug"`
	WebhookProvider string `gorm:"type:varchar(255);column:webhook_provider" json:"webhook_provider"`
	WebhookUrl      string `gorm:"type:text;column:webhook_url" json:"webhook_url"`
}

/*
	for filtering field use like this for [carts] table:
	- carts.quantity -> even for current table filtering, always call the table name like this
	- products.name -> filter using products table with field name ->
	remember to not using struct field -> always use real tables and field name
*/

func (Project) TableName() string {
	return "projects"
}
