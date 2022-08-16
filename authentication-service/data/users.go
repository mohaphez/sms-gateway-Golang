package data

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name,omitempty"`
	UserName  string    `json:"username,omitempty"`
	Password  string    `json:"-"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) UserIsActive(username string) (int32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	qstring := `select status from users where email = $1 and status=1`
	var isActive int32
	err := db.QueryRow(ctx, qstring, username).Scan(&isActive)

	switch {
	case err == sql.ErrNoRows:
		log.Printf("no user with username %s", username)
		return 0, err
	case err != nil:
		log.Println(err)
		return 0, err
	default:
		log.Printf("username is %s\n", username)
		return isActive, nil
	}
}

// Create new user and return the ID of newly inserted row

func (u *User) Insert(user User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)

	if err != nil {
		return 0, err
	}

	var userID int
	query := `insert into users (email, name, username, password, status, created_at, updated_at)
	values ($1, $2, $3, $4, $5, $6, $7) returning id`

	err = db.QueryRow(ctx, query,
		user.Email,
		user.Name,
		hashedPassword,
		user.Status,
		time.Now(),
		time.Now()).Scan(&userID)

	if err != nil {
		return 0, err
	}

	return userID, nil

}

// ResetPassword is the method we use to change a user password.
func (u *User) ResetPassword(password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `update users set password = $1 where id = $2`
	_, err = db.Exec(ctx, stmt, hashedPassword, u.ID)
	if err != nil {
		return err
	}

	return nil
}

// PasswordMatch is the method for check user password matche with database hashed string password.
func (u *User) PasswordMatch(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))

	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}
