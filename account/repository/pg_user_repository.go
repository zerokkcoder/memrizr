package repository

import (
	"context"
	"log"
	"memrizr/model"
	"memrizr/model/apperrors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// PGUserRepository 用户存储层实现
type PGUserRepository struct {
	DB *sqlx.DB
}

// NewUserRepository 实例化PGUserRepository
func NewUserRepository(db *sqlx.DB) model.UserRepository {
	return &PGUserRepository{
		DB: db,
	}
}

// FindByID 通过 ID 查找用户
func (r *PGUserRepository) FindByID(ctx context.Context, uid uuid.UUID) (*model.User, error) {
	user := &model.User{}

	query := "SELECT * FROM users WHERE uid=$1"

	if err := r.DB.Get(user, query, uid); err != nil {
		return user, apperrors.NewNotFound("uid", uid.String())
	}

	return user, nil
}

// Create 创建用户
func (r *PGUserRepository) Create(ctx context.Context, u *model.User) error {
	query := "INSERT INTO users (email, password) VALUES ($1, $2) RETURN *"

	if err := r.DB.Get(u, query, u.Email, u.Password); err != nil {
		// 检验 唯一
		if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
			log.Printf("Could not create a user with email: %v. Reason: %v\n", u.Email, err.Code.Name())
			return apperrors.NewConflict("email", u.Email)
		}

		log.Printf("Could not create a user with email: %v. Reason: %v\n", u.Email, err.Code.Name())
		return apperrors.NewInternal()
	}

	return nil
}
