package votes

import (
	"log"

	"github.com/hemanta212/hackernews-go-graphql/internal/links"
	database "github.com/hemanta212/hackernews-go-graphql/internal/pkg/db/mysql"
	"github.com/hemanta212/hackernews-go-graphql/internal/users"
)

type Vote struct {
	ID      string
	Link    *links.Link
	VotedBy *users.User
}

func (vote Vote) Save() int64 {
	stmt, err := database.Db.Prepare("INSERT INTO Votes(LinkID, UserID) VALUES(?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(vote.Link.ID, vote.VotedBy.ID)
	if err != nil {
		log.Fatal(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal("Error: ", err.Error())
	}
	return id
}
