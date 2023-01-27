package links

import (
	"database/sql"
	"log"
	"strconv"

	database "github.com/hemanta212/hackernews-go-graphql/internal/pkg/db/mysql"
	"github.com/hemanta212/hackernews-go-graphql/internal/users"
)

type Link struct {
	ID          string
	Description string
	Url         string
	PostedBy    *users.User
}

func (link Link) Save() int64 {
	stmt, err := database.Db.Prepare("INSERT INTO Links(Description, Url, UserID) VALUES(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(link.Description, link.Url, link.PostedBy.ID)
	if err != nil {
		log.Fatal(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal("Error: ", err.Error())
	}
	return id
}

func GetLinkByID(id int) (*Link, error) {
	stmt, err := database.Db.Prepare(
		"SELECT L.Description, L.Url, L.UserID, U.Username, U.Email from Links L inner join Users U on L.UserID = U.ID WHERE L.ID=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	link := &Link{ID: strconv.Itoa(id)}
	user := &users.User{}
	err = stmt.QueryRow(id).Scan(&link.Description, &link.Url, &user.ID, &user.Username, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		} else {
			log.Fatal(err)
		}
	}
	link.PostedBy = user

	return link, nil
}

func GetAll() []Link {
	stmt, err := database.Db.Prepare(
		"SELECT L.ID, L.Description, L.Url, L.UserID, U.Username, U.Email from Links L inner join Users U on L.UserID = U.ID")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var links []Link
	for rows.Next() {
		var link Link
		user := &users.User{}
		err := rows.Scan(&link.ID, &link.Description, &link.Url, &user.ID, &user.Username, &user.Email)
		if err != nil {
			log.Fatal(err)
		}
		link.PostedBy = user
		links = append(links, link)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return links
}
