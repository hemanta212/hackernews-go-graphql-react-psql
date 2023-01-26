package users

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"

	database "github.com/hemanta212/hackernews-go-graphql/internal/pkg/db/mysql"
)

type User struct {
	ID       string  `json:"id"`
	Username string  `json:"name"`
	Password string  `json:"password"`
	Email    *string `json:"email"`
}

func (user *User) Save() (int64, error) {
	stmt, err := database.Db.Prepare("INSERT INTO Users(Username, Password, Email) VALUES(?, ?, ?)")
	if err != nil {
		return -1, err
	}
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return -1, err
	}
	res, err := stmt.Exec(user.Username, hashedPassword, user.Email)
	if err != nil {
		return -1, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (user *User) Authenticate(rawPassword string) bool {
	return checkPasswordHash(user.Password, rawPassword)
}

func GetUserByUsername(userName string) (*User, error) {
	stmt, err := database.Db.Prepare("SELECT password, ID, email FROM Users WHERE username=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var user User
	err = stmt.QueryRow(userName).Scan(&user.Password, &user.ID, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		} else {
			log.Fatal(err)
		}
	}
	return &user, nil
}

func GetUserIdByUsername(username string) (int64, error) {
	stmt, err := database.Db.Prepare("SELECT id from Users where username=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var userID int64
	err = stmt.QueryRow(username).Scan(&userID)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return 0, err
	}
	return userID, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
