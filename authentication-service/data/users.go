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
	Password  string    `json:"password"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func createUserstable() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	qstring := `CREATE TABLE IF NOT EXISTS users (
		id serial PRIMARY KEY,
		username VARCHAR ( 50 ) UNIQUE NOT NULL,
		name VARCHAR ( 50 ) NULL,
		password VARCHAR ( 255 ) NOT NULL,
		email VARCHAR ( 255 ) UNIQUE NOT NULL,
		status INT  default 1,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL
	);`

	_, err := db.Exec(ctx, qstring)

	if err != nil {
		log.Println("Database can't create users table")
		log.Println(err)
		return err
	}

	return nil
}

func (u *User) UserIsActive(username string) ([2]int, string, error) {
	var userinfo [2]int
	var userpass string
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	qstring := `select id,status,password from users where username = $1 and status = 1`
	err := db.QueryRow(ctx, qstring, username).Scan(&userinfo[0], &userinfo[1], &userpass)

	switch {
	case err == sql.ErrNoRows:
		log.Printf("no user with username %s", username)
		return [2]int{0, 0}, "-", err
	case err != nil:
		log.Println(err)
		return [2]int{0, 0}, "-", err
	default:
		log.Printf("username is %s\n", username)
		return userinfo, userpass, nil
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
		user.UserName,
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
func (u *User) PasswordMatch(password string, hashpassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashpassword), []byte(password))

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
