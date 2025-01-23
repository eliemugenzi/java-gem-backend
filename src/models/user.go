package models

import "time"

/*
* User model
 */
type User struct {
	ID string `json:"id" gorm:"primary_key"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	Email string `json:"email"`
	Password string `json:"password"`
	Role     string `gorm:"type:enum('ADMIN','USER');default:'USER';not null"` // Role field
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
