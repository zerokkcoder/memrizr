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
type pgUserRepository struct {
	DB *sqlx.DB
}

// NewUserRepository 实例化PGUserRepository
func NewUserRepository(db *sqlx.DB) model.UserRepository {
	return &pgUserRepository{
		DB: db,
	}
}

// FindByID 通过 ID 查找用户
func (r *pgUserRepository) FindByID(ctx context.Context, uid uuid.UUID) (*model.User, error) {
	user := &model.User{}

	query := "SELECT * FROM users WHERE uid=$1;"

	if err := r.DB.GetContext(ctx, user, query, uid); err != nil {
		return user, apperrors.NewNotFound("uid", uid.String())
	}

	return user, nil
}

// FindByEmail 通过 Email 查找用户
func (r *pgUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}

	query := "SELECT * FROM users WHERE email=$1;"

	if err := r.DB.GetContext(ctx, user, query, email); err != nil {
		log.Printf("Unable to get user with email address: %v. Err: %v\n", email, err)
		return user, apperrors.NewNotFound("email", email)
	}

	return user, nil
}

// Create 创建用户
func (r *pgUserRepository) Create(ctx context.Context, u *model.User) error {
	query := "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING *;"

	if err := r.DB.GetContext(ctx, u, query, u.Email, u.Password); err != nil {
		// 检验 唯一
		if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
			log.Printf("Could not create a user with email: %v. Reason: %v\n", u.Email, err.Code.Name())
			return apperrors.NewConflict("email", u.Email)
		}

		log.Printf("Could not create a user with email: %v. Reason: %v\n", u.Email, err)
		return apperrors.NewInternal()
	}

	return nil
}

// Update 更新用户
func (r *pgUserRepository) Update(ctx context.Context, u *model.User) error {
	query := `
		UPDATE users
		SET name=:name, email=:email, website=:website
		WHERE uid=:uid
		RETURNING *;
	`
	nsmt, err := r.DB.PrepareNamedContext(ctx, query)
	if err != nil {
		log.Printf("Unable to prepare user update query: %v\n", err)
		return apperrors.NewInternal()
	}

	if err := nsmt.GetContext(ctx, u, u); err != nil {
		log.Printf("Unable to update details for user: %v\n", u)
		return apperrors.NewInternal()
	}

	return nil
}
