package repository

import (
	"context"

	"github.com/oguzhantasimaz/Go-Clean-Architecture-Template/domain"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetUsers(ctx context.Context) ([]*domain.User, error)
	GetUserById(ctx context.Context, id int) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, userId int) error
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetUsers(ctx context.Context) ([]*domain.User, error) {
	var users []*domain.User
	err := r.db.SelectContext(ctx, &users, "SELECT * FROM users")
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) GetUserById(ctx context.Context, id int) (*domain.User, error) {
	user := domain.User{}
	err := r.db.GetContext(ctx, &user, `SELECT * FROM users WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := domain.User{}
	err := r.db.GetContext(ctx, &user, `SELECT * FROM users WHERE email = $1`, email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Commit()

	if user.GoogleId != "" {
		var id int
		err := tx.QueryRowContext(ctx,
			`INSERT INTO users (email, google_id, name, profile_picture) VALUES ($1, $2, $3, $4) RETURNING id`,
			user.Email, user.GoogleId, user.Name, user.ProfilePicture,
		).Scan(&id)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		user.Id = id
		return user, nil
	}

	var id int
	err = tx.QueryRowContext(ctx,
		`INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`,
		user.Email, user.Password,
	).Scan(&id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	user.Id = id

	return user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Commit()

	if user.Password != "" {
		encryptedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(user.Password),
			bcrypt.DefaultCost,
		)
		if err != nil {
			return err
		}
		user.Password = string(encryptedPassword)
	}

	fieldsQuery := ""
	if user.Email != "" {
		fieldsQuery += "email = :email,"
	}
	if user.Name != "" {
		fieldsQuery += "name = :name,"
	}
	if user.Password != "" {
		fieldsQuery += "password = :password,"
	}
	if user.Phone != "" {
		fieldsQuery += "phone = :phone,"
	}
	fieldsQuery = fieldsQuery[:len(fieldsQuery)-1]

	_, err = tx.NamedExecContext(ctx, "UPDATE users SET "+fieldsQuery+" WHERE id = :id", user)
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *userRepository) DeleteUser(ctx context.Context, userId int) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Commit()

	_, err = tx.ExecContext(ctx, "DELETE FROM users WHERE id = $1", userId)
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
