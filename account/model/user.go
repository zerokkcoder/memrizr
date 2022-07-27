package model

import (
	"github.com/google/uuid"
)

// User 用户模型
type User struct {
	UID      uuid.UUID `db:"uid" json:"uid"`
	Email    string    `db: "email" json:"email"`
	Password string    `db: "password" json:"-"`
	Name     string    `db: "name" json:"name"`
	ImageURL string    `db: "image_url" json:"image_url"`
	Website  string    `db: "website" json:"website"`
}
