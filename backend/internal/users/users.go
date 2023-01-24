package users

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"

	database "github.com/hemanta212/hackernews-go-graphql/internal/pkg/db/mysql"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"name"`
	Password string `json:"password"`
}

func (user *User) Save() {
	stmt, err := database.Db.Prepare("INSERT INTO Users(Username, Password) VALUES(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(user.Username, hashedPassword)
	if err != nil {
		log.Fatal(err)
	}

}

func (user *User) Authenticate() bool {
	stmt, err := database.Db.Prepare("SELECT password FROM Users WHERE username=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var hashedPassword string
	err = stmt.QueryRow(user.Username).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			log.Fatal(err)
		}
	}
	return checkPasswordHash(user.Password, hashedPassword)
}

func GetUserIdByUsername(username string) (int, error) {
	stmt, err := database.Db.Prepare("SELECT id from Users where username=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var userID int
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
