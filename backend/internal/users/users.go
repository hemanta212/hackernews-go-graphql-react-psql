package users

import (
	"database/sql"
	"log"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	database "github.com/hemanta212/hackernews-go-graphql/internal/pkg/db/postgresql"
)

type User struct {
	ID       string  `json:"id"`
	Username string  `json:"name"`
	Password string  `json:"password"`
	Email    *string `json:"email"`
}

func (user *User) Save() (int, error) {
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return -1, err
	}

	var lastInsertId int
	err = database.Db.QueryRow("INSERT INTO Users(Username, Password, Email) VALUES($1,$2,$3) returning id;",
		user.Username, hashedPassword, user.Email,
	).Scan(&lastInsertId)

	user.ID = strconv.Itoa(lastInsertId)

	if err != nil {
		return -1, err
	}

	return lastInsertId, nil
}

func (user *User) Authenticate(rawPassword string) bool {
	return checkPasswordHash(rawPassword, user.Password)
}

func GetUserByUsername(userName string) (*User, error) {
	stmt, err := database.Db.Prepare("SELECT ID, password, email FROM Users WHERE username=$1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	user := &User{Username: userName}
	err = stmt.QueryRow(userName).Scan(&user.ID, &user.Password, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		} else {
			log.Fatal(err)
		}
	}
	return user, nil
}

func GetUserByID(id int) (*User, error) {
	stmt, err := database.Db.Prepare("SELECT username, password, email FROM Users WHERE id=$1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	user := &User{ID: strconv.Itoa(id)}
	err = stmt.QueryRow(id).Scan(&user.Username, &user.Password, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		} else {
			log.Fatal(err)
		}
	}
	return user, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
