package store

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plaintText *string
	hash       []byte
}

// ******************** password.Set ********************
func (p *password) Set(plainTextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintText = &plainTextPassword
	p.hash = hash
	return nil
}

// ******************** password.Matches ********************
func (p *password) Matches(plainTextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainTextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err // internal server error
		}
	}
	return true, nil
}

type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"user_name"`
	Email        string    `json:"email"`
	PasswordHash password  `json:"-"`
	Bio          string    `json:"bio"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

var AnonymousUser = &User{}

// ******************** User.IsAnonymous ********************
func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

type PostgresUserStore struct {
	db *sql.DB
}

type UserStore interface {
	CreateUser(*User) error
	GetUserByUsername(username string) (*User, error)
	UpdateUser(*User) error
	GetUserToken(scope, tokenPlainText string) (*User, error)
}

// ******************** constructor: NewPostgresUserStore ********************
func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{
		db: db,
	}
}

// ******************** PostgresUserStore.CreateUser ********************
func (s *PostgresUserStore) CreateUser(user *User) error {
	query := `
		INSERT INTO users (user_name, email, password_hash, bio)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	err := s.db.QueryRow(query, user.Username, user.Email, user.PasswordHash.hash, user.Bio).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

// ******************** PostgresUserStore.GetUserByUsername ********************
func (s *PostgresUserStore) GetUserByUsername(username string) (*User, error) {
	user := &User{
		PasswordHash: password{},
	}

	query := `
		SELECT id, user_name, email, password_hash, bio, created_at, updated_at
		FROM users
		WHERE user_name = $1
	`
	err := s.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash.hash,
		&user.Bio,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	fmt.Println("user:", user)
	return user, nil
}

// ******************** PostgresUserStore.UpdateUser ********************
func (s *PostgresUserStore) UpdateUser(user *User) error {
	query := `
		UPDATE users
		set username = $1, email = $2, bio = $3, updated_at = current_timestamp
		WHERE id = $4
		RETURNING updated_at
	`
	result, err := s.db.Exec(query, user.Username, user.Email, user.Bio, user.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// ******************** PostgresUserStore.GetUserToken ********************
func (s *PostgresUserStore) GetUserToken(scope string, plainTextPassword string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(plainTextPassword))

	query := `
		SELECT u.id, u.user_name, u.email, u.password_hash, u.bio, u.created_at, u.updated_at
		FROM users u
		INNER JOIN tokens t ON t.user_id = u.id
		WHERE t.hash = $1 AND t.scope = $2 and t.expiry > $3
	`
	user := &User{
		PasswordHash: password{},
	}

	err := s.db.QueryRow(query, tokenHash[:], scope, time.Now()).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash.hash,
		&user.Bio,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return user, nil
}
